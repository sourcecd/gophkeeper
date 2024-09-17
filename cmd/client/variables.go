package main

import (
	"flag"

	"github.com/sourcecd/gophkeeper/internal/options"
)

// load configuration from cmd line
// TODO from env too
func loadConfiguration(opt *options.ClientOptions) {
	clientFlags(opt)
}

// client cmdline flags parse
func clientFlags(opt *options.ClientOptions) {
	flag.StringVar(&opt.GrpcAddr, "grpc-addr", "localhost:2135", "grpc server address")
	flag.StringVar(&opt.HttpAddr, "http-addr", "localhost:8080", "listen http server address")
	flag.StringVar(&opt.CAfile, "ca-file-path", "", "if empty, uses embeded cert")
	flag.Parse()
}
