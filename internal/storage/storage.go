package storage

import keeperproto "github.com/sourcecd/gophkeeper/proto"

type ServerStorage interface {
	RegisterUser()
	AuthUser()
	SyncPut()
	SyncGet()
}

type ClientStorage interface {
	SyncPut()
	SyncGet(protodata *[]*keeperproto.Data) error
	PutItem(name, vType string, value []byte) error
	GetItem(name string, valueType *string, value *[]byte) error
	DelItem(name string) error
}
