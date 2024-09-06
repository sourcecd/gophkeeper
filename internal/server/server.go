package server

import (
	"context"
	"log"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/storage"
)

func Run(ctx context.Context, opt *options.ServerOptions) {
	_, err := storage.PgBaseInit(ctx, opt.Dsn)
	if err != nil {
		log.Fatal(err)
	}
}