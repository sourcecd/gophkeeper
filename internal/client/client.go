package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type handlers struct {
	store storage.ClientStorage
}

func syncPush(ctx context.Context, conn *grpc.ClientConn) error {
	c := keeperproto.NewSyncClient(conn)
	resp, err := c.Push(ctx, &keeperproto.SyncPushRequest{
		Data: []*keeperproto.Data{
			{
				Name:    "TEST2",
				Type:    keeperproto.Data_Type(keeperproto.Data_Type_value["TEXT"]),
				Payload: []byte("OK"),
			},
		},
	})
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func syncPull(ctx context.Context, conn *grpc.ClientConn) error {
	c := keeperproto.NewSyncClient(conn)
	resp, err := c.Pull(ctx, &keeperproto.SyncPullRequest{
		Name: []string{},
	})
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func grpcConn(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func newHandlers(s storage.ClientStorage) *handlers {
	return &handlers{store: s}
}

func (h *handlers) postItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("can't read body")
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		if err := h.store.PutItem(chi.URLParam(r, "name"), chi.URLParam(r, "type"), req); err != nil {
			slog.Error("can't store value")
			http.Error(w, "can't store value", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("STORED"))
	}
}

func (h *handlers) getItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *handlers) delItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func chiRouter(h *handlers) chi.Router {
	r := chi.NewRouter()

	r.Post("/add/{type}/{name}", h.postItem())
	r.Get("/get/{name}", h.getItem())
	r.Delete("/del/{name}", h.delItem())

	return r
}

func Run(ctx context.Context, opt *options.ClientOptions) {
	conn, err := grpcConn(opt.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	err = syncPush(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	err = syncPull(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	h := newHandlers(storage.NewInMemory())
	log.Printf("Starting http server: %s", opt.HttpAddr)
	if err := http.ListenAndServe(opt.HttpAddr, chiRouter(h)); err != nil {
		log.Fatal(err)
	}
}
