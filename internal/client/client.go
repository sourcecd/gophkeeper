package client

import (
	"context"
	"fmt"
	"log"

	"github.com/sourcecd/gophkeeper/internal/options"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
}
