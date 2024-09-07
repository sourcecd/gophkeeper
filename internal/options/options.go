package options

type ServerOptions struct {
	Dsn      string
	GrpcAddr string
}

type ClientOptions struct {
	GrpcAddr string
}
