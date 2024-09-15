package main

import (
	"context"

	"github.com/sourcecd/gophkeeper/internal/client"
	"github.com/sourcecd/gophkeeper/internal/options"
)

func main() {
	// Print Build args
	printBuildFlags()

	ctx := context.Background()
	var opt options.ClientOptions

	loadConfiguration(&opt)

	client.Run(ctx, &opt)
}
