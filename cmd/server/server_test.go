package main

import (
	"testing"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestServerVars(t *testing.T) {
	opt := &options.ServerOptions{}

	//run load cfg
	loadConfiguration(opt)

	assert.Equal(t, "dbname=gophkeeper", opt.Dsn)
	assert.Equal(t, "localhost:2135", opt.GrpcAddr)
	assert.Equal(t, "key", opt.SecurityKeyFile)
	assert.Equal(t, "", opt.CertPemFile)
	assert.Equal(t, "", opt.KeyPemFile)
}
