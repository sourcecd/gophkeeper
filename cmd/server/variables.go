package main

import (
	"flag"

	"github.com/sourcecd/gophkeeper/internal/options"
)

// load configuration from cmd line
// TODO from env too
func loadConfiguration(opt *options.ServerOptions) {
	serverFlags(opt)
}

// server cmdline flags parse
func serverFlags(opt *options.ServerOptions) {
	flag.StringVar(&opt.Dsn, "dsn", "dbname=gophkeeper", "dsn for postgres")
	flag.StringVar(&opt.GrpcAddr, "grpc-addr", "localhost:2135", "listen grpc server address")
	flag.StringVar(&opt.SecurityKeyFile, "sec-key-file", "", "if empty, uses embeded security key for crypt")
	flag.StringVar(&opt.CertPemFile, "cert-pem-file", "", "if empty, uses embeded cert")
	flag.StringVar(&opt.KeyPemFile, "key-pem-file", "", "if empty, uses embeded key")
	flag.Parse()
}
