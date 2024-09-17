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

// For testing only
//
//go:embed certs/ca.crt
var embedCerts embed.FS

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
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(500000000)))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
