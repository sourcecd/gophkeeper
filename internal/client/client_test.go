package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestBase(t *testing.T) {
	var err error
	ctx := context.Background()
	sClient := NewSyncClient(ctx, conn, inmemory)
	h := newHandlers(sClient)

	r := chi.NewRouter()
	r.Post("/add/{type}/{name}/{description}", h.postItem())
	require.NoError(t, err)

	srv := httptest.NewServer(r)
	defer srv.Close()

	_, err = http.Post(srv.URL+"/add/text/ok/ok", "", nil)
	require.NoError(t, err)
}
