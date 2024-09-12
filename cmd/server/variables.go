package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/sourcecd/gophkeeper/internal/options"
)

func loadConfiguration(opt *options.ServerOptions) {
	serverFlags(opt)
	f, err := os.Open(opt.SecurityKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	opt.SecurityKey = string(b)
}

func serverFlags(opt *options.ServerOptions) {
	flag.StringVar(&opt.Dsn, "dsn", "dbname=gophkeeper", "dsn for postgres")
	flag.StringVar(&opt.GrpcAddr, "grpc-addr", "localhost:2135", "listen grpc server address")
	flag.StringVar(&opt.SecurityKeyFile, "sec-key-file", "key", "security key for crypt")
	flag.Parse()
}
