package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"github.com/stretchr/testify/require"
)

const testToken = "XXX.YYY.ZZZ"

type SyncTest struct {
	inmem *storage.InMemoryStore
}

func (sy *SyncTest) SyncPush(token string, proto []*keeperproto.Data) error {
	return nil
}
func (sy *SyncTest) SyncPull(token string) error {
	return nil
}
func (sy *SyncTest) AuthUser(login, password string, token *string) error {
	*token = testToken
	return nil
}
func (sy *SyncTest) RegisterUser(login, password string, token *string) error {
	*token = testToken
	return nil
}
func (sy *SyncTest) Store() storage.ClientStorage {
	return sy.inmem
}

func TestClient(t *testing.T) {
	sClient := &SyncTest{
		inmem: storage.NewInMemory(),
	}
	cli := &http.Client{}

	h := newHandlers(sClient)

	srv := httptest.NewServer(chiRouter(h))
	defer srv.Close()

	testCases := []struct {
		name       string
		url        string
		body       io.Reader
		result     string
		method     string
		statusCode int
	}{
		{
			name:       "postItem",
			url:        "/add/text/test/desc",
			body:       strings.NewReader("test1"),
			result:     "STORED\n",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
		{
			name:       "getItem",
			url:        "/get/test",
			body:       nil,
			result:     "test1",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
		{
			name:       "listallItem",
			url:        "/listall",
			body:       nil,
			result:     "Name | Type | Description\n--------------------------\ntest | TEXT | desc\n\n",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
		{
			name:       "delItem",
			url:        "/del/test",
			body:       nil,
			result:     "REMOVED\n",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
		},
		{
			name:       "registerUser",
			url:        "/register/testuser/testpass",
			body:       nil,
			result:     testToken + "\n",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
		{
			name:       "authUser",
			url:        "/auth/testuser/testpass",
			body:       nil,
			result:     testToken + "\n",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			req, err := http.NewRequest(v.method, srv.URL+v.url, v.body)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+testToken)
			resp, err := cli.Do(req)
			require.NoError(t, err)
			require.Equal(t, v.statusCode, resp.StatusCode)

			b, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			require.NoError(t, err)
			require.Equal(t, v.result, string(b))
		})
	}
}
