// Package options for client and server
package options

// ServerOptions options for server
type ServerOptions struct {
	Dsn             string
	GrpcAddr        string
	SecurityKeyFile string
	SecurityKey     string
}

// ClientOptions options for client
type ClientOptions struct {
	GrpcAddr string
	HttpAddr string
}
