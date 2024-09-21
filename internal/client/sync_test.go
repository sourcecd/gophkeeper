package client

import (
	"context"
	"testing"

	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.nhat.io/grpcmock"
)

func TestSyncItems(t *testing.T) {
	t.Parallel()

	_, d := grpcmock.MockServerWithBufConn(
		grpcmock.RegisterService(keeperproto.RegisterSyncServer),
		func(s *grpcmock.Server) {
			s.ExpectUnary("gophkeeper.Sync/Push").Once().
				WithPayload(&keeperproto.SyncPushRequest{Data: []*keeperproto.Data{
					{
						Name: "TEST",
					},
				}}).
				Return(&keeperproto.SyncPushResponse{Error: "OK"})
		},
		func(s *grpcmock.Server) {
			s.ExpectUnary("gophkeeper.Sync/Pull").Once().
				WithPayload(&keeperproto.SyncPullRequest{Name: []string{}}).
				Return(&keeperproto.SyncPullResponse{Data: []*keeperproto.Data{{Name: "OK"}}})
		},
		func(s *grpcmock.Server) {
			s.ExpectUnary("gophkeeper.Sync/RegisterUser").Once().
				WithPayload(&keeperproto.AuthRequest{Login: "login", Password: "password"}).
				Return(&keeperproto.AuthResponse{Token: "test"})
		},
		func(s *grpcmock.Server) {
			s.ExpectUnary("gophkeeper.Sync/AuthUser").Once().
				WithPayload(&keeperproto.AuthRequest{Login: "login", Password: "password"}).
				Return(&keeperproto.AuthResponse{Token: "test"})
		},
	)(t)

	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(d), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	syncli := NewSyncClient(ctx, conn, storage.NewInMemory())

	// Push
	err = syncli.SyncPush("qwe", []*keeperproto.Data{
		{
			Name: "TEST",
		},
	})
	require.NoError(t, err)

	// Pull
	err = syncli.SyncPull("qwe")
	require.NoError(t, err)

	// RegisterUser
	var testToken string
	err = syncli.RegisterUser("login", "password", &testToken)
	require.NoError(t, err)

	// Auth User
	err = syncli.AuthUser("login", "password", &testToken)
	require.NoError(t, err)
}

func TestGrpcConn(t *testing.T) {
	_, err := grpcConn("", "")
	require.NoError(t, err)
}
