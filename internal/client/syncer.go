package client

import (
	"context"
	"crypto/x509"
	"embed"
	"log"
	"os"

	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	"github.com/sourcecd/gophkeeper/internal/storage"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const maxRecvMsgSize = 524288000

// For testing only
//
//go:embed certs/ca.crt
var embedCerts embed.FS

// SyncClientIface interface for client syncer
type SyncClientIface interface {
	SyncPush(token string, proto []*keeperproto.Data) error
	SyncPull(token string) error
	AuthUser(login, password string, token *string) error
	RegisterUser(login, password string, token *string) error
	Store() storage.ClientStorage
}

// SyncClient struct
type SyncClient struct {
	ctx   context.Context
	conn  *grpc.ClientConn
	store storage.ClientStorage
}

// NewSyncClient sync client init
func NewSyncClient(ctx context.Context, conn *grpc.ClientConn, store storage.ClientStorage) *SyncClient {
	return &SyncClient{
		ctx:   ctx,
		conn:  conn,
		store: store,
	}
}

// Store return link to client storage iface
func (sy *SyncClient) Store() storage.ClientStorage {
	return sy.store
}

// RegisterUser send register user grpc request to server
func (sy *SyncClient) RegisterUser(login, password string, token *string) error {
	c := keeperproto.NewSyncClient(sy.conn)

	resp, err := c.RegisterUser(sy.ctx, &keeperproto.AuthRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}
	*token = resp.Token
	return nil
}

// AuthUser send auth user grpc request to server
func (sy *SyncClient) AuthUser(login, password string, token *string) error {
	c := keeperproto.NewSyncClient(sy.conn)

	resp, err := c.AuthUser(sy.ctx, &keeperproto.AuthRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}
	*token = resp.Token
	return nil
}

// SyncPush push data to server by grpc
func (sy *SyncClient) SyncPush(token string, data []*keeperproto.Data) error {
	c := keeperproto.NewSyncClient(sy.conn)

	ctx := metadata.NewOutgoingContext(sy.ctx, metadata.MD{
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

// SyncPull pull data from server by grpc
func (sy *SyncClient) SyncPull(token string) error {
	c := keeperproto.NewSyncClient(sy.conn)

	ctx := metadata.NewOutgoingContext(sy.ctx, metadata.MD{
		"token": {token},
	})

	resp, err := c.Pull(ctx, &keeperproto.SyncPullRequest{
		Name: []string{},
	})
	if err != nil {
		return err
	}

	if err := sy.store.SyncPut(resp.Data); err != nil {
		return err
	}
	log.Printf("Synced records from server: %d", len(resp.Data))

	return nil
}

// use tls certificate for grpc client
func generateTLSCreds(caPath string) (credentials.TransportCredentials, error) {
	var (
		certb []byte
		err   error
	)

	if caPath == "" {
		certb, err = embedCerts.ReadFile("certs/ca.crt")
	} else {
		certb, err = os.ReadFile(caPath)
	}
	if err != nil {
		return nil, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(certb) {
		return nil, fixederrors.ErrCertificateLoad
	}

	return credentials.NewClientTLSFromCert(cp, ""), nil
}

// grpc connection to server
func grpcConn(addr, caFile string) (*grpc.ClientConn, error) {
	tlsCreds, err := generateTLSCreds(caFile)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(tlsCreds),
		// TODO make stream
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRecvMsgSize)))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
