package storage

import (
	"sync"

	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
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

func (c *ClientInMemory) SyncPut(protodata []*keeperproto.Data) error {
	c.Lock()
	defer c.Unlock()
	for _, v := range protodata {
		c.data[v.Name] = valueStore{
			valueType: v.Type.String(),
			value:     v.Payload,
		}
	}
	return nil
}

func (c *ClientInMemory) SyncGet(protodata *[]*keeperproto.Data) error {
	c.RLock()
	defer c.RUnlock()
	for k, v := range c.data {
		*protodata = append(*protodata, &keeperproto.Data{
			Name:    k,
			Type:    keeperproto.Data_Type(keeperproto.Data_Type_value[v.valueType]),
			Payload: v.value,
		})
	}
	return nil
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

func (c *ClientInMemory) DelItem(name string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[name]; ok {
		delete(c.data, name)
		return nil
	}
	return fixederrors.ErrNoValue
}
