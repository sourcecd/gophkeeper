package main

import (
	"flag"

	"github.com/sourcecd/gophkeeper/internal/options"
)

func loadConfiguration(opt *options.ServerOptions) {
	serverFlags(opt)
}

func serverFlags(opt *options.ServerOptions) {
	flag.StringVar(&opt.Dsn, "dsn", "dbname=gophkeeper", "dsn for postgres")
	flag.StringVar(&opt.GrpcAddr, "grpc-addr", "localhost:2135", "grpc server address")
	flag.Parse()
}
