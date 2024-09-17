package main

import (
	"testing"

	"github.com/sourcecd/gophkeeper/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestClientVars(t *testing.T) {
	opt := &options.ClientOptions{}

	//run load cfg
	loadConfiguration(opt)

	assert.Equal(t, "localhost:2135", opt.GrpcAddr)
	assert.Equal(t, "localhost:8080", opt.HttpAddr)
	assert.Equal(t, "", opt.CAfile)
}
