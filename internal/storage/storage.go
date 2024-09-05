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
	PutItem()
	GetItem()
}
