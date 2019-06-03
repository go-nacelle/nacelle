package grpc

type Config struct {
	GRPCHost string `env:"grpc_host" file:"grpc_host" default:"0.0.0.0"`
	GRPCPort int    `env:"grpc_port" file:"grpc_port" default:"6000"`
}
