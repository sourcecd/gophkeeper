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
	"github.com/sourcecd/gophkeeper/internal/dataparser"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
)

type handlers struct {
	store storage.ClientStorage
	ctx   context.Context
	conn  *grpc.ClientConn
}

func checkTypeValue(r *http.Request) error {
	s := chi.URLParam(r, "type")
	if _, ok := keeperproto.Data_DType_value[strings.ToUpper(s)]; !ok {
		return fixederrors.ErrUnkType
	}
	return nil
}

func baseTokenCheck(r *http.Request, token *string) error {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer") {
		if s := strings.Split(authHeader, " "); len(s) == 2 {
			*token = s[1]
			return nil
		}
	}
	return fixederrors.ErrInvalidToken
}

func newHandlers(ctx context.Context, s storage.ClientStorage, conn *grpc.ClientConn) *handlers {
	return &handlers{
		store: s,
		ctx:   ctx,
		conn:  conn,
	}
}

func (h *handlers) postItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		if err := baseTokenCheck(r, &token); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
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
		name := chi.URLParam(r, "name")
		dtype := strings.ToUpper(chi.URLParam(r, "type"))
		desc := chi.URLParam(r, "description")

		parsedData, err := dataparser.Dataparser(dtype, req).Parse()
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := syncPush(h.ctx, h.conn, token, nil, []*keeperproto.Data{
			{
				Name:        name,
				Optype:      keeperproto.Data_OpType(keeperproto.Data_OpType_value["ADD"]),
				Dtype:       keeperproto.Data_DType(keeperproto.Data_DType_value[dtype]),
				Payload:     parsedData,
				Description: desc,
			},
		}); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.store.PutItem(name, dtype, req, desc); err != nil {
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

func (h *handlers) delItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		if err := baseTokenCheck(r, &token); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		n := chi.URLParam(r, "name")
		if err := syncPush(h.ctx, h.conn, token, nil, []*keeperproto.Data{
			{
				Name:   n,
				Optype: keeperproto.Data_OpType(keeperproto.Data_OpType_value["DELETE"]),
			},
		}); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.store.DelItem(n); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("REMOVED\n"))
	}
}

func (h *handlers) registerUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		login := chi.URLParam(r, "login")
		password := chi.URLParam(r, "password")

		if err := registerUser(h.ctx, h.conn, login, password, &token); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token + "\n"))
	}
}

func (h *handlers) authUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		login := chi.URLParam(r, "login")
		password := chi.URLParam(r, "password")

		if err := authUser(h.ctx, h.conn, login, password, &token); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := syncPull(h.ctx, h.conn, token, h.store); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token + "\n"))
	}
}

func chiRouter(h *handlers) chi.Router {
	r := chi.NewRouter()

	r.Post("/add/{type}/{name}/{description}", h.postItem())
	r.Get("/get/{name}", h.getItem())
	r.Delete("/del/{name}", h.delItem())
	r.Post("/register/{login}/{password}", h.registerUser())
	r.Post("/auth/{login}/{password}", h.authUser())

	return r
}

func Run(ctx context.Context, opt *options.ClientOptions) {
	inmemory := storage.NewInMemory()
	conn, err := grpcConn(opt.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	h := newHandlers(ctx, inmemory, conn)
	log.Printf("Starting http server: %s", opt.HttpAddr)
	if err := http.ListenAndServe(opt.HttpAddr, chiRouter(h)); err != nil {
		log.Fatal(err)
	}
}
