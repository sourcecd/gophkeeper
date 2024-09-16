package client

import (
	"context"
	"log"

	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// send register user grpc request to server
func registerUser(ctx context.Context, conn *grpc.ClientConn, login, password string, token *string) error {
	c := keeperproto.NewSyncClient(conn)

	resp, err := c.RegisterUser(ctx, &keeperproto.AuthRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}
	*token = resp.Token
	return nil
}

// send auth user grpc request to server
func authUser(ctx context.Context, conn *grpc.ClientConn, login, password string, token *string) error {
	c := keeperproto.NewSyncClient(conn)

	resp, err := c.AuthUser(ctx, &keeperproto.AuthRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}
	*token = resp.Token
	return nil
}

// push data to server by grpc
func syncPush(ctx context.Context, conn *grpc.ClientConn, token string, data []*keeperproto.Data) error {
	c := keeperproto.NewSyncClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"token": {token},
	})

	resp, err := c.Push(ctx, &keeperproto.SyncPushRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	log.Printf("Synced records to server: %d", len(data))

	if resp.Error != "" {
		log.Println(resp.Error)
	}
	return nil
}

// pull data from server by grpc
func syncPull(ctx context.Context, conn *grpc.ClientConn, token string, store storage.ClientStorage) error {
	c := keeperproto.NewSyncClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"token": {token},
	})

	resp, err := c.Pull(ctx, &keeperproto.SyncPullRequest{
		Name: []string{},
	})
	if err != nil {
		return err
	}

	if err := store.SyncPut(resp.Data); err != nil {
		return err
	}
	log.Printf("Synced records from server: %d", len(resp.Data))

	return nil
}

// grpc connection to server
func grpcConn(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		// TODO make stream
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(500000000)))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
