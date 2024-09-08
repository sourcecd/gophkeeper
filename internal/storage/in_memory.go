package storage

import (
	"sync"

	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
)

type valueStore struct {
	valueType string
	value     []byte
}

type ClientInMemory struct {
	sync.RWMutex
	data map[string]valueStore
}

func NewInMemory() *ClientInMemory {
	return &ClientInMemory{
		data: make(map[string]valueStore),
	}
}

func (c *ClientInMemory) SyncPut() {
}

func (c *ClientInMemory) SyncGet() {
}

func (c *ClientInMemory) PutItem(name, vType string, value []byte) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[name]; ok {
		return fixederrors.ErrAlreadyExists
	}
	c.data[name] = valueStore{
		valueType: vType,
		value:     value,
	}
	return nil
}

func (c *ClientInMemory) GetItem(name string, valueType *string, value *[]byte) error {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.data[name]; ok {
		*value = v.value
		*valueType = v.valueType
		return nil
	}
	return fixederrors.ErrNoValue
}
