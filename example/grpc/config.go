package grpc

type (
	Config struct {
		GRPCPort int `env:"GRPC_PORT" default:"6000"`
	}

	configToken struct{}
)

var ConfigToken = configToken{}
