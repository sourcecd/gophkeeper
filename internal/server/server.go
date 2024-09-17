// Package server with grpc transport for client
package server

import (
	"context"
	"crypto/tls"
	"embed"
	"log"
	"net"

	"github.com/sourcecd/gophkeeper/internal/auth"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const maxRecvMsgSize = 524288000

// For testing only
//
//go:embed certs/server.crt certs/server.key
var embedCerts embed.FS

// SyncServer main grpc protocol struct
type SyncServer struct {
	keeperproto.UnimplementedSyncServer
	store storage.ServerStorage
	jwtm  *auth.JWTManager
}

// check and unpack jwt token
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

// NewSyncServer create server instance
func NewSyncServer(db storage.ServerStorage, jwtm *auth.JWTManager) *SyncServer {
	return &SyncServer{store: db, jwtm: jwtm}
}

// RegisterUser server grpc method for register user
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

// AuthUser server grpc method for authorize user
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

// Push data to server storage
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

// Pull data from server storage
func (s *SyncServer) Pull(ctx context.Context, in *keeperproto.SyncPullRequest) (*keeperproto.SyncPullResponse, error) {
	var (
		data   []*keeperproto.Data
		userid int64
	)
	if err := s.checkToken(ctx, &userid); err != nil {
		return nil, err
	}
	if err := s.store.SyncGet(ctx, in.Name, &data, userid); err != nil {
		return nil, err
	}
	log.Printf("Synced records to client: %d", len(data))
	return &keeperproto.SyncPullResponse{
		Data: data,
	}, nil
}

// use tls certificate for grpc server
func generateTLSCreds() (credentials.TransportCredentials, error) {
	certFile := "certs/server.crt"
	keyFile := "certs/server.key"

	certb, err := embedCerts.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyb, err := embedCerts.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(certb, keyb)
	if err != nil {
		return nil, err
	}

	return credentials.NewServerTLSFromCert(&tlsCert), nil
}

// grpc server configure
func startGrpcServer(addr string, syncServer *SyncServer) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsCreds, err := generateTLSCreds()
	if err != nil {
		return err
	}

	// TODO make stream
	s := grpc.NewServer(grpc.MaxRecvMsgSize(maxRecvMsgSize), grpc.Creds(tlsCreds))
	keeperproto.RegisterSyncServer(s, syncServer)
	log.Printf("Starting grpc server on: %s", addr)
	if err := s.Serve(l); err != nil {
		return err
	}
	return nil
}

// Run server
func Run(ctx context.Context, opt *options.ServerOptions) {
	db, err := storage.PgBaseInit(ctx, opt.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := startGrpcServer(opt.GrpcAddr, NewSyncServer(db, auth.NewJWTManager(opt.SecurityKey))); err != nil {
		log.Fatal(err)
	}
}
