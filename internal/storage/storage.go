// Package storage with interfaces and implementations
package storage

import (
	"context"

	"github.com/sourcecd/gophkeeper/internal/auth"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

// ServerStorage interface
type ServerStorage interface {
	RegisterUser(ctx context.Context, reg *auth.User, userid *int64) error
	AuthUser(ctx context.Context, reg *auth.User, userid *int64) error
	SyncPut(ctx context.Context, data []*keeperproto.Data, userid int64) error
	SyncGet(ctx context.Context, names []string, data *[]*keeperproto.Data, userid int64) error
}

// ClientStorage interface
type ClientStorage interface {
	SyncPut(protodata []*keeperproto.Data) error
	SyncGet(protodata *[]*keeperproto.Data) error
	PutItem(name, vType string, value []byte, desc string) error
	GetItem(name string, valueType *string, value *[]byte) error
	DelItem(name string) error
	ListItems(items *[]ListItems) error
}
