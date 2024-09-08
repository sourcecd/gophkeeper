package client

import (
	"context"
	"fmt"

	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func syncPush(ctx context.Context, conn *grpc.ClientConn, store storage.ClientStorage) error {
	c := keeperproto.NewSyncClient(conn)
	var data []*keeperproto.Data

	if err := store.SyncGet(&data); err != nil {
		return err
	}

	resp, err := c.Push(ctx, &keeperproto.SyncPushRequest{
		Data: data,
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
