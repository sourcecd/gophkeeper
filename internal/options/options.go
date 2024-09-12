package options

type ServerOptions struct {
	Dsn             string
	GrpcAddr        string
	SecurityKeyFile string
	SecurityKey     string
}

type ClientOptions struct {
	GrpcAddr string
	HttpAddr string
}
