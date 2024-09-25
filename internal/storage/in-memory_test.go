package storage

import (
	"testing"

	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	var (
		valType string
		value   []byte
		valData []*keeperproto.Data
		items   []ListItems
	)
	inmemory := NewInMemory()

	err := inmemory.SyncPut([]*keeperproto.Data{
		{
			Name:        "test1",
			Dtype:       keeperproto.Data_DType(keeperproto.Data_DType_value["TEXT"]),
			Optype:      keeperproto.Data_OpType(keeperproto.Data_OpType_value["ADD"]),
			Payload:     []byte{200},
			Description: "OK",
		},
	})
	require.NoError(t, err)

	err = inmemory.GetItem("test1", &valType, &value)
	require.NoError(t, err)
	require.Equal(t, value, []byte{200})

	err = inmemory.SyncGet(&valData)
	require.NoError(t, err)
	require.Equal(t, "test1", valData[0].Name)

	err = inmemory.ListItems(&items)
	require.NoError(t, err)
	require.Equal(t, "test1", items[0].Name)

	err = inmemory.PutItem("test2", "TEXT", []byte{102}, "OK")
	require.NoError(t, err)

	err = inmemory.DelItem("test1")
	require.NoError(t, err)
}
