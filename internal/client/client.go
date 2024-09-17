// Package client for cache and send items
package client

import (
	"context"
	"fmt"
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

// handlers type with methods for REST api
type handlers struct {
	store storage.ClientStorage
	ctx   context.Context
	conn  *grpc.ClientConn
}

// validate type of data
func checkTypeValue(r *http.Request) error {
	s := chi.URLParam(r, "type")
	if _, ok := keeperproto.Data_DType_value[strings.ToUpper(s)]; !ok {
		return fixederrors.ErrUnkType
	}
	return nil
}

// check authorization header
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

// format items to string
func itemsStringView(items []storage.ListItems) string {
	s := "Name | Type | Description\n--------------------------\n"
	for _, v := range items {
		s += fmt.Sprintf("%s | %s | %s\n", v.Name, v.DType, v.Desc)
	}
	return s
}

// create handler instance
func newHandlers(ctx context.Context, s storage.ClientStorage, conn *grpc.ClientConn) *handlers {
	return &handlers{
		store: s,
		ctx:   ctx,
		conn:  conn,
	}
}

// accept request to client from user and send it to server
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

		if err := syncPush(h.ctx, h.conn, token, []*keeperproto.Data{
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

// fetch item from client by user request
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

// del item by user request from client and sync delete with server
func (h *handlers) delItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		if err := baseTokenCheck(r, &token); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		n := chi.URLParam(r, "name")
		if err := syncPush(h.ctx, h.conn, token, []*keeperproto.Data{
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

// send register user request to server and get jwt token response
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

// send auth user request to server and get jwt token response
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

// list all items from client by user request
func (h *handlers) listAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listItems []storage.ListItems
		if err := h.store.ListItems(&listItems); err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s := itemsStringView(listItems)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s + "\n"))
	}
}

// handler router
func chiRouter(h *handlers) chi.Router {
	r := chi.NewRouter()

	r.Post("/add/{type}/{name}/{description}", h.postItem())
	r.Get("/get/{name}", h.getItem())
	r.Delete("/del/{name}", h.delItem())
	r.Post("/register/{login}/{password}", h.registerUser())
	r.Post("/auth/{login}/{password}", h.authUser())
	r.Get("/listall", h.listAll())

	return r
}

// Run client daemon
func Run(ctx context.Context, opt *options.ClientOptions) {
	inmemory := storage.NewInMemory()
	conn, err := grpcConn(opt.GrpcAddr, opt.CAfile)
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
