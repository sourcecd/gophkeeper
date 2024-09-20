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

func TestGetItems(t *testing.T) {
	t.Parallel()

	expected := &keeperproto.SyncPushResponse{Error: "OK"}

	_, d := grpcmock.MockServerWithBufConn(
		grpcmock.RegisterService(keeperproto.RegisterSyncServer),
		func(s *grpcmock.Server) {
			s.ExpectUnary("gophkeeper.Sync/Push").
				WithPayload(&keeperproto.SyncPushRequest{Data: []*keeperproto.Data{
					{
						Name: "TEST",
					},
				}}).
				Return(expected)
		},
	)(t)

	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(d), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	syncli := NewSyncClient(ctx, conn, storage.NewInMemory())
	err = syncli.SyncPush("qwe", []*keeperproto.Data{
		{
			Name: "TEST",
		},
	})
	require.NoError(t, err)
}
