package grpc

import (
	"fmt"
)

type (
	Config struct {
		GRPCPort int `env:"grpc_port" default:"6000"`
	}

	configToken string
)

var ConfigToken = NewConfigToken("default")

func NewConfigToken(name string) interface{} {
	return configToken(fmt.Sprintf("nacelle-base-grpc-%s", name))
}
