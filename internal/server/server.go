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
	store storage.ServerStorage
}

func NewSyncServer(db storage.ServerStorage) *SyncServer {
	return &SyncServer{store: db}
}

func (s *SyncServer) Push(ctx context.Context, in *keeperproto.SyncPushRequest) (*keeperproto.SyncPushResponse, error) {
	if len(in.Data) == 1 && in.Data[0].Optype == keeperproto.Data_OpType(keeperproto.Data_OpType_value["DELETE"]) {

	}
	if err := s.store.SyncPut(ctx, in.Data); err != nil {
		return &keeperproto.SyncPushResponse{
			Error: err.Error(),
		}, nil
	}
	log.Printf("Synced records from client: %d", len(in.Data))
	return &keeperproto.SyncPushResponse{
		Error: "",
	}, nil
}

func (s *SyncServer) Pull(ctx context.Context, in *keeperproto.SyncPullRequest) (*keeperproto.SyncPullResponse, error) {
	var data []*keeperproto.Data
	if err := s.store.SyncGet(ctx, in.Name, &data); err != nil {
		return nil, err
	}
	log.Printf("Synced records to client: %d", len(data))
	return &keeperproto.SyncPullResponse{
		Data: data,
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
