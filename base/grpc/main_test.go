package grpc

//go:generate go-mockgen github.com/efritz/nacelle/config -i Config -o mock_config_test.go -f

import (
	"net"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/config"
	"github.com/efritz/nacelle/service"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ServerSuite{})
	})
}

//
// Config

type emptyConfig struct{}

func makeConfig(base *Config) nacelle.Config {
	cfg := NewMockConfig()
	cfg.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...config.TagModifier) error {
		c := target.(*Config)
		c.GRPCPort = base.GRPCPort
		return nil
	})

	return cfg
}

//
//  Injection

type A struct{ X int }
type B struct{ X float64 }

func makeBadContainer() nacelle.ServiceContainer {
	container, _ := service.NewContainer()
	container.Set("A", &B{})
	return container
}

//
// Server Helpers

func getDynamicPort(listener net.Listener) int {
	return listener.Addr().(*net.TCPAddr).Port
}
