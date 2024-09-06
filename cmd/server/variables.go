package main

import (
	"flag"

	"github.com/sourcecd/gophkeeper/internal/options"
)

func loadConfiguration(opt *options.ServerOptions) {
	serverFlags(opt)
}

func serverFlags(opt *options.ServerOptions) {
	flag.StringVar(&opt.Dsn, "dsn", "", "dsn for postgres")
	flag.Parse()
}