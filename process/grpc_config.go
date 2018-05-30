package process

import (
	"github.com/efritz/nacelle"
)

type (
	GRPCConfig struct {
		GRPCPort int `env:"grpc_port" default:"6000"`
	}

	grpcConfigToken string
)

var GRPCConfigToken = grpcConfigToken("nacelle-process-grpc")

// RegisterGRPCConfigs adds the required configs for a GRPC server to the given map. If any tag
// modifiers are supplied, they are run over each of the required configs (this may require
// some knowledge about package internals).
func RegisterGRPCConfigs(m map[interface{}]interface{}, modifiers ...nacelle.TagModifier) error {
	c, err := nacelle.ApplyTagModifiers(&GRPCConfig{}, modifiers...)
	if err != nil {
		return err
	}

	m[GRPCConfigToken] = c
	return nil
}
