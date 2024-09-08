package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
)

type SyncServer struct {
	keeperproto.UnimplementedSyncServer
	store storage.ServerStorage
}

func NewSyncServer(db storage.ServerStorage) *SyncServer {
	return &SyncServer{store: db}
}

func (s *SyncServer) Push(ctx context.Context, in *keeperproto.SyncPushRequest) (*keeperproto.SyncPushResponse, error) {
	if err := s.store.SyncPut(ctx, in.Data); err != nil {
		return &keeperproto.SyncPushResponse{
			Error: err.Error(),
		}, nil
	}
	return &keeperproto.SyncPushResponse{
		Error: "",
	}, nil
}

func (s *SyncServer) Pull(ctx context.Context, in *keeperproto.SyncPullRequest) (*keeperproto.SyncPullResponse, error) {
	fmt.Println(in)
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

func startGrpcServer(addr string, syncServer *SyncServer) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	keeperproto.RegisterSyncServer(s, syncServer)
	log.Printf("Starting grpc server on: %s", addr)
	if err := s.Serve(l); err != nil {
		return err
	}
	return nil
}

func Run(ctx context.Context, opt *options.ServerOptions) {
	db, err := storage.PgBaseInit(ctx, opt.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := startGrpcServer(opt.GrpcAddr, NewSyncServer(db)); err != nil {
		log.Fatal(err)
	}
}
