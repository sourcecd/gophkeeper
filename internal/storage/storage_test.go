package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sourcecd/gophkeeper/internal/auth"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
	"github.com/stretchr/testify/require"
)

const (
	pass     = "mypass"
	hashpass = "$2a$10$WKank/GVLIhij4iWIohBb.EpwaFVHX4xFexBfuLrZ.tDk20OLurRe"
)

var (
	mock sqlmock.Sqlmock
	db   *sql.DB
	err  error
	pgdb *PgDB
)

func TestCreatePGDB(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	pgdb = &PgDB{
		db: db,
	}

	// prepare test
	mock.ExpectPrepare(putDataRequest)
	mock.ExpectPrepare(selectAllDataRequest)
	mock.ExpectPrepare(deleteItemRequest)
	mock.ExpectPrepare(createUserRequest)
	mock.ExpectPrepare(getUserRequest)

	err = pgdb.PrepStmt()
	require.NoError(t, err)
}

func TestSyncPut(t *testing.T) {
	userid := int64(10)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(putDataRequest).WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(10, "test1", 1, []byte{101}, "test1")
	mock.ExpectCommit()
	mock.ExpectExec(deleteItemRequest).WillReturnResult(sqlmock.NewResult(1, 1))

	// base sync put
	err = pgdb.SyncPut(ctx, []*keeperproto.Data{{
		Name:        "test1",
		Dtype:       keeperproto.Data_DType(keeperproto.Data_DType_value["TEXT"]),
		Optype:      keeperproto.Data_OpType(keeperproto.Data_OpType_value["ADD"]),
		Payload:     []byte{101},
		Description: "test1",
	}}, userid)
	require.NoError(t, err)

	// delete item
	err = pgdb.SyncPut(ctx, []*keeperproto.Data{{
		Name:        "test1",
		Dtype:       keeperproto.Data_DType(keeperproto.Data_DType_value["TEXT"]),
		Optype:      keeperproto.Data_OpType(keeperproto.Data_OpType_value["DELETE"]),
		Payload:     []byte{101},
		Description: "test1",
	}}, userid)
	require.NoError(t, err)
}

func TestSyncGet(t *testing.T) {
	userid := int64(10)
	ctx := context.Background()

	mock.ExpectQuery(selectAllDataRequest).WithArgs(10).WillReturnRows(sqlmock.NewRows([]string{"name", "type", "payload", "description"}).AddRow("test1", "1", []byte{}, "test1"))

	err = pgdb.SyncGet(ctx, []string{}, &[]*keeperproto.Data{{}}, userid)
	require.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	var userid int64
	ctx := context.Background()

	mock.ExpectQuery(createUserRequest).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	err = pgdb.RegisterUser(ctx, &auth.User{}, &userid)
	require.NoError(t, err)
}

func TestAuthUser(t *testing.T) {
	var userid int64
	ctx := context.Background()

	mock.ExpectQuery(getUserRequest).WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).AddRow(10, "test1", hashpass))

	err = pgdb.AuthUser(ctx, &auth.User{
		Username:       "test1",
		HashedPassword: pass,
	}, &userid)
	require.NoError(t, err)
}
