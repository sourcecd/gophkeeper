// Package main for client
package main

import (
	"context"

	"github.com/sourcecd/gophkeeper/internal/client"
	"github.com/sourcecd/gophkeeper/internal/options"
)

func main() {
	// print Build args
	printBuildFlags()

	ctx := context.Background()
	var opt options.ClientOptions

	// load client configuration from cmdline
	// TODO from env too
	loadConfiguration(&opt)

	// main client start
	client.Run(ctx, &opt)
}
