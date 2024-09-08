package storage

import (
	"context"
	"database/sql"
	"embed"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	keeperproto "github.com/sourcecd/gophkeeper/proto"
)

const (
	putDataRequest = "INSERT INTO data (name, type, payload) VALUES ($1, $2, $3)"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PgDB struct {
	db                 *sql.DB
	stmtPutDataRequest *sql.Stmt
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

func (pg *PgDB) SyncPut(ctx context.Context, data []*keeperproto.Data) error {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range data {
		if _, err := tx.StmtContext(ctx, pg.stmtPutDataRequest).ExecContext(ctx, v.Name, v.Type, v.Payload); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (pg *PgDB) SyncGet() error {
	return nil
}

func (pg *PgDB) RegisterUser() error {
	return nil
}

func (pg *PgDB) AuthUser() error {
	return nil
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
