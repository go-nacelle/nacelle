package grpc

type Config struct {
	GRPCPort int `env:"grpc_port" default:"6000"`
}
