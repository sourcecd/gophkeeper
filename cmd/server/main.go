// Package main for server
package main

import (
	"context"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/sourcecd/gophkeeper/internal/server"
)

func main() {
	// print Build args
	printBuildFlags()

	ctx := context.Background()
	var opt options.ServerOptions

	// load server configuration from cmdline
	// TODO from env too
	loadConfiguration(&opt)

	// main server start
	server.Run(ctx, &opt)
}
