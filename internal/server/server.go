package server

import (
	"context"
	"log"
	"net"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
)

type SyncServer struct {
	keeperproto.UnimplementedSyncServer
}

func (s *SyncServer) Push(ctx context.Context, in *keeperproto.SyncPushRequest) (*keeperproto.SyncPushResponse, error) {
	return &keeperproto.SyncPushResponse{
		Error: "NoError",
	}, nil
}

func (s *SyncServer) Pull(ctx context.Context, in *keeperproto.SyncPullRequest) (*keeperproto.SyncPullResponse, error) {
	return &keeperproto.SyncPullResponse{
		Data: []*keeperproto.Data{
			{
				Name:    "TEST",
				Type:    keeperproto.Data_Type(keeperproto.Data_Type_value["TEXT"]),
				Payload: []byte("TEST"),
			},
		},
	}, nil
}

func startGrpcServer(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	keeperproto.RegisterSyncServer(s, &SyncServer{})
	if err := s.Serve(l); err != nil {
		return err
	}
	return nil
}

func Run(ctx context.Context, opt *options.ServerOptions) {
	_, err := storage.PgBaseInit(ctx, opt.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := startGrpcServer(opt.GrpcAddr); err != nil {
		log.Fatal(err)
	}
}
