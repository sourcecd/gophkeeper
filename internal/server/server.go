package server

import (
	"context"
	"log"
	"net"

	"github.com/sourcecd/gophkeeper/internal/auth"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type SyncServer struct {
	keeperproto.UnimplementedSyncServer
	store storage.ServerStorage
	jwtm  *auth.JWTManager
}

func (s *SyncServer) checkToken(ctx context.Context, userid *int64) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fixederrors.ErrInvalidToken
	}
	mdToken := md["token"]
	if len(mdToken) == 0 {
		return fixederrors.ErrInvalidToken
	}
	clm, err := s.jwtm.Verify(mdToken[0])
	if err != nil {
		return err
	}
	*userid = clm.UserID
	return nil
}

func NewSyncServer(db storage.ServerStorage, jwtm *auth.JWTManager) *SyncServer {
	return &SyncServer{store: db, jwtm: jwtm}
}

// TODO move jwt to ???
func (s *SyncServer) RegisterUser(ctx context.Context, in *keeperproto.AuthRequest) (*keeperproto.AuthResponse, error) {
	var userid int64
	user, err := s.jwtm.PrepareUser(in.Login, in.Password)
	if err != nil {
		return nil, err
	}
	if err := s.store.RegisterUser(ctx, user, &userid); err != nil {
		return nil, err
	}
	token, err := s.jwtm.Generate(userid)
	if err != nil {
		return nil, err
	}
	return &keeperproto.AuthResponse{
		Token: token,
	}, nil
}

// TODO move jwt to ???
func (s *SyncServer) AuthUser(ctx context.Context, in *keeperproto.AuthRequest) (*keeperproto.AuthResponse, error) {
	var userid int64
	if err := s.store.AuthUser(ctx, &auth.User{
		Username: in.Login,
		// not hashed
		HashedPassword: in.Password,
	}, &userid); err != nil {
		return nil, err
	}
	token, err := s.jwtm.Generate(userid)
	if err != nil {
		return nil, err
	}
	return &keeperproto.AuthResponse{
		Token: token,
	}, nil
}

func (s *SyncServer) Push(ctx context.Context, in *keeperproto.SyncPushRequest) (*keeperproto.SyncPushResponse, error) {
	var userid int64
	if err := s.checkToken(ctx, &userid); err != nil {
		return nil, err
	}
	if err := s.store.SyncPut(ctx, in.Data, userid); err != nil {
		return &keeperproto.SyncPushResponse{
			Error: err.Error(),
		}, err
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

	if err := startGrpcServer(opt.GrpcAddr, NewSyncServer(db, auth.NewJWTManager(opt.SecurityKey))); err != nil {
		log.Fatal(err)
	}
}
