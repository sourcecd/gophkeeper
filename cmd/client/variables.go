package main

import (
	"flag"

	"github.com/sourcecd/gophkeeper/internal/options"
)

func loadConfiguration(opt *options.ClientOptions) {
	ClientFlags(opt)
}

func ClientFlags(opt *options.ClientOptions) {
	flag.StringVar(&opt.GrpcAddr, "grpc-addr", "localhost:2135", "grpc server address")
	flag.StringVar(&opt.HttpAddr, "http-addr", "localhost:8080", "listen http server address")
	flag.Parse()
}
