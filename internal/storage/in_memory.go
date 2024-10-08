package storage

import (
	"sync"

	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

// struct for save value in-memory
type valueStore struct {
	valueType string
	value     []byte
	desc      string
}

// ListItems struct for metainfo
type ListItems struct {
	Name  string
	DType string
	Desc  string
}

// InMemoryStore in-memory client storage
type InMemoryStore struct {
	sync.RWMutex
	data map[string]valueStore
}

// NewInMemory init client in-memory storage
func NewInMemory() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]valueStore),
	}
}

// SyncPut put data (with server sync) to client in-memory storage
func (c *InMemoryStore) SyncPut(protodata []*keeperproto.Data) error {
	c.Lock()
	defer c.Unlock()
	for _, v := range protodata {
		c.data[v.Name] = valueStore{
			valueType: v.Dtype.String(),
			value:     v.Payload,
			desc:      v.Description,
		}
	}
	return nil
}

// SyncGet get data (with server sync) from client in-memory storage
// Not used, for potencial full sync from client to server
func (c *InMemoryStore) SyncGet(protodata *[]*keeperproto.Data) error {
	c.RLock()
	defer c.RUnlock()
	for k, v := range c.data {
		*protodata = append(*protodata, &keeperproto.Data{
			Name:        k,
			Dtype:       keeperproto.Data_DType(keeperproto.Data_DType_value[v.valueType]),
			Payload:     v.value,
			Description: v.desc,
		})
	}
	return nil
}

// PutItem to in-memory client storage (one item)
func (c *InMemoryStore) PutItem(name, vType string, value []byte, desc string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[name]; ok {
		return fixederrors.ErrAlreadyExists
	}
	c.data[name] = valueStore{
		valueType: vType,
		value:     value,
		desc:      desc,
	}
	return nil
}

// GetItem from in-memory client storage (one item)
func (c *InMemoryStore) GetItem(name string, valueType *string, value *[]byte) error {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.data[name]; ok {
		*value = v.value
		*valueType = v.valueType
		return nil
	}
	return fixederrors.ErrNoValue
}

// ListItems metainfo from in-memory client storage
func (c *InMemoryStore) ListItems(items *[]ListItems) error {
	c.RLock()
	defer c.RUnlock()
	for k, v := range c.data {
		*items = append(*items, ListItems{
			Name:  k,
			DType: v.valueType,
			Desc:  v.desc,
		})
	}
	return nil
}

// DelItem from in-memory client storage
func (c *InMemoryStore) DelItem(name string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[name]; ok {
		delete(c.data, name)
		return nil
	}
	return fixederrors.ErrNoValue
}
