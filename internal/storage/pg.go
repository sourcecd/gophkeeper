package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sourcecd/gophkeeper/internal/auth"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

const (
	putDataRequest       = "INSERT INTO data (id, name, type, payload) VALUES ($1, $2, $3, $4)"
	selectAllDataRequest = "SELECT name, type, payload FROM data WHERE id = $1"
	deleteItemRequest    = "DELETE FROM data WHERE id = $1 AND name = $2"
	createUserRequest    = "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id"
	getUserRequest       = "SELECT id, login, password FROM users WHERE login = $1"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PgDB struct {
	db                       *sql.DB
	stmtPutDataRequest       *sql.Stmt
	stmtSelectAllDataRequest *sql.Stmt
	stmtDeleteItemRequest    *sql.Stmt
	stmtCreateUserRequest    *sql.Stmt
	stmtGetUserRequest       *sql.Stmt
}

func NewPgDB(dsn string) (*PgDB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &PgDB{
		db: db,
	}, nil
}

func (pg *PgDB) PrepStmt() error {
	var err error
	pg.stmtPutDataRequest, err = pg.db.Prepare(putDataRequest)
	if err != nil {
		return err
	}
	pg.stmtSelectAllDataRequest, err = pg.db.Prepare(selectAllDataRequest)
	if err != nil {
		return err
	}
	pg.stmtDeleteItemRequest, err = pg.db.Prepare(deleteItemRequest)
	if err != nil {
		return err
	}
	pg.stmtCreateUserRequest, err = pg.db.Prepare(createUserRequest)
	if err != nil {
		return err
	}
	pg.stmtGetUserRequest, err = pg.db.Prepare(getUserRequest)
	if err != nil {
		return err
	}
	log.Println("requests prepared")
	return nil
}

func (pg *PgDB) CreateDatabaseScheme(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, pg.db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (pg *PgDB) SyncPut(ctx context.Context, data []*keeperproto.Data, userid int64) error {
	if len(data) == 1 && data[0].Optype == keeperproto.Data_OpType(keeperproto.Data_OpType_value["DELETE"]) {
		res, err := pg.stmtDeleteItemRequest.ExecContext(ctx, userid, data[0].Name)
		if err != nil {
			return err
		}
		r, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if r == 0 {
			return fixederrors.ErrRecordNotFound
		}
		return nil
	}
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range data {
		if _, err := tx.StmtContext(ctx, pg.stmtPutDataRequest).ExecContext(ctx, userid, v.Name, v.Dtype, v.Payload); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return fixederrors.ErrRecordAlreadyExists
			}
			return err
		}
	}
	return tx.Commit()
}

func (pg *PgDB) SyncGet(ctx context.Context, names []string, data *[]*keeperproto.Data, userid int64) error {
	var (
		name    string
		vType   string
		payload []byte
	)
	if len(names) == 0 {
		rows, err := pg.stmtSelectAllDataRequest.QueryContext(ctx, userid)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&name, &vType, &payload); err != nil {
				return err
			}
			*data = append(*data, &keeperproto.Data{
				Name:    name,
				Dtype:   keeperproto.Data_DType(keeperproto.Data_DType_value[vType]),
				Payload: payload,
			})
		}
		if err := rows.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (pg *PgDB) RegisterUser(ctx context.Context, reg *auth.User, userid *int64) error {
	err := pg.stmtCreateUserRequest.QueryRowContext(ctx, reg.Username, reg.HashedPassword).Scan(userid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fixederrors.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (pg *PgDB) AuthUser(ctx context.Context, reg *auth.User, userid *int64) error {
	var (
		login,
		password string
	)
	row := pg.stmtGetUserRequest.QueryRowContext(ctx, reg.Username)
	if err := row.Scan(userid, &login, &password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fixederrors.ErrUserNotExists
		}
		return err
	}
	if reg.IsCorrectPassword(password) {
		return nil
	}
	userid = nil
	return fixederrors.ErrUserNotExists
}

func PgBaseInit(ctx context.Context, dsn string) (*PgDB, error) {
	db, err := NewPgDB(dsn)
	if err != nil {
		return nil, err
	}
	if err := db.CreateDatabaseScheme(ctx); err != nil {
		return nil, err
	}
	if err := db.PrepStmt(); err != nil {
		return nil, err
	}
	return db, nil
}
