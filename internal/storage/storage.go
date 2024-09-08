package storage

import (
	"context"

	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

type ServerStorage interface {
	RegisterUser() error
	AuthUser() error
	SyncPut(ctx context.Context, data []*keeperproto.Data) error
	SyncGet() error
}

type ClientStorage interface {
	SyncPut()
	SyncGet(protodata *[]*keeperproto.Data) error
	PutItem(name, vType string, value []byte) error
	GetItem(name string, valueType *string, value *[]byte) error
	DelItem(name string) error
}
