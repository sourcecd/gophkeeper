package storage

import (
	"sync"

	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

type valueStore struct {
	valueType string
	value     []byte
	desc      string
}

type ListItems struct {
	Name  string
	DType string
	Desc  string
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
			valueType: v.Dtype.String(),
			value:     v.Payload,
			desc:      v.Description,
		}
	}
	return nil
}

func (c *ClientInMemory) SyncGet(protodata *[]*keeperproto.Data) error {
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

func (c *ClientInMemory) PutItem(name, vType string, value []byte, desc string) error {
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

func (c *ClientInMemory) ListItems(items *[]ListItems) error {
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

func (c *ClientInMemory) DelItem(name string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[name]; ok {
		delete(c.data, name)
		return nil
	}
	return fixederrors.ErrNoValue
}
