package main

import (
	"context"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/server"
)

func main() {
	ctx := context.Background()
	var opt options.ServerOptions

	loadConfiguration(&opt)

	server.Run(ctx, &opt)
}
