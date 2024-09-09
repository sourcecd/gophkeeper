package storage

import (
	"context"

	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

type ServerStorage interface {
	RegisterUser() error
	AuthUser() error
	SyncPut(ctx context.Context, data []*keeperproto.Data) error
	SyncGet(ctx context.Context, names []string, data *[]*keeperproto.Data) error
}

type ClientStorage interface {
	SyncPut(protodata []*keeperproto.Data) error
	SyncGet(protodata *[]*keeperproto.Data) error
	PutItem(name, vType string, value []byte) error
	GetItem(name string, valueType *string, value *[]byte) error
	DelItem(name string) error
}
