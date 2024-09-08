package storage

type ServerStorage interface {
	RegisterUser()
	AuthUser()
	SyncPut()
	SyncGet()
}

type ClientStorage interface {
	SyncPut()
	SyncGet()
	PutItem(name, vType string, value []byte) error
	GetItem(name string, valueType *string, value *[]byte) error
	DelItem(name string) error
}
