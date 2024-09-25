package server

import (
	"context"
	"testing"

	"github.com/sourcecd/gophkeeper/internal/auth"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

// const secKey = "Nee5chahNi2eepao"

type JWTtest struct {
}

func (j *JWTtest) PrepareUser(username string, password string) (*auth.User, error) {
	return nil, nil
}
func (j *JWTtest) Generate(userid int64) (string, error) {
	return "", nil
}
func (j *JWTtest) Verify(accessToken string) (*auth.UserClaims, error) {
	return &auth.UserClaims{UserID: 10}, nil
}

type ServerTest struct {
}

func (s *ServerTest) RegisterUser(ctx context.Context, reg *auth.User, userid *int64) error {
	return nil
}
func (s *ServerTest) AuthUser(ctx context.Context, reg *auth.User, userid *int64) error {
	return nil
}
func (s *ServerTest) SyncPut(ctx context.Context, data []*keeperproto.Data, userid int64) error {
	return nil
}
func (s *ServerTest) SyncGet(ctx context.Context, names []string, data *[]*keeperproto.Data, userid int64) error {
	return nil
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{"token": []string{"testok"}})
	srv := NewSyncServer(&ServerTest{}, &JWTtest{})

	// Register user
	_, err := srv.RegisterUser(ctx, &keeperproto.AuthRequest{})
	require.NoError(t, err)

	// Auth user
	_, err = srv.AuthUser(ctx, &keeperproto.AuthRequest{})
	require.NoError(t, err)

	// Sync Pull
	_, err = srv.Pull(ctx, &keeperproto.SyncPullRequest{})
	require.NoError(t, err)

	// Sync Push
	_, err = srv.Push(ctx, &keeperproto.SyncPushRequest{})
	require.NoError(t, err)
}

func TestCheckTLS(t *testing.T) {
	tls, err := generateTLSCreds("", "")
	require.NoError(t, err)
	require.NotNil(t, tls)
}
