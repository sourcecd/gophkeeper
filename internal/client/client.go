package client

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"

	// "time"

	"github.com/go-chi/chi/v5"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
)

type handlers struct {
	store storage.ClientStorage
	ctx   context.Context
}

func checkTypeValue(r *http.Request) error {
	s := chi.URLParam(r, "type")
	if _, ok := keeperproto.Data_DType_value[strings.ToUpper(s)]; !ok {
		return fixederrors.ErrUnkType
	}
	return nil
}

func newHandlers(ctx context.Context, s storage.ClientStorage) *handlers {
	return &handlers{
		store: s,
		ctx:   ctx,
	}
}

func (h *handlers) postItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkTypeValue(r); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("can't read body")
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		if err := h.store.PutItem(chi.URLParam(r, "name"), strings.ToUpper(chi.URLParam(r, "type")), req); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("STORED\n"))
	}
}

func (h *handlers) getItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			valType string
			val     []byte
		)
		if err := h.store.GetItem(chi.URLParam(r, "name"), &valType, &val); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(val)
	}
}

func (h *handlers) delItem(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := chi.URLParam(r, "name")
		if err := h.store.DelItem(n); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := syncPush(h.ctx, conn, nil, []*keeperproto.Data{
			{
				Name:   n,
				Optype: keeperproto.Data_OpType(keeperproto.Data_OpType_value["DELETE"]),
			},
		}); err != nil {
			slog.Error(err.Error())
			http.Error(w, "error when delete data from server", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("REMOVED\n"))
	}
}

func chiRouter(h *handlers, conn *grpc.ClientConn) chi.Router {
	r := chi.NewRouter()

	r.Post("/add/{type}/{name}", h.postItem())
	r.Get("/get/{name}", h.getItem())
	r.Delete("/del/{name}", h.delItem(conn))

	return r
}

func Run(ctx context.Context, opt *options.ClientOptions) {
	inmemory := storage.NewInMemory()
	conn, err := grpcConn(opt.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = syncPull(ctx, conn, inmemory)
	if err != nil {
		log.Fatal(err)
	}

	h := newHandlers(ctx, inmemory)
	log.Printf("Starting http server: %s", opt.HttpAddr)
	if err := http.ListenAndServe(opt.HttpAddr, chiRouter(h, conn)); err != nil {
		log.Fatal(err)
	}
}
