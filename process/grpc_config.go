package process

import (
	"fmt"
)

type (
	GRPCConfig struct {
		GRPCPort int `env:"grpc_port" default:"6000"`
	}

	grpcConfigToken string
)

var GRPCConfigToken = MakeGRPCConfigToken("default")

func MakeGRPCConfigToken(name string) interface{} {
	return grpcConfigToken(fmt.Sprintf("nacelle-process-grpc-%s", name))
}
