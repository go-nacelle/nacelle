package http

//go:generate go-mockgen github.com/go-nacelle/nacelle/config -i Config -o mock_config_test.go -f

import (
	"net"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/go-nacelle/nacelle"
	"github.com/go-nacelle/nacelle/config"
	"github.com/go-nacelle/nacelle/service"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ConfigSuite{})
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
		c.HTTPPort = base.HTTPPort
		c.HTTPCertFile = base.HTTPCertFile
		c.HTTPKeyFile = base.HTTPKeyFile
		c.RawShutdownTimeout = base.RawShutdownTimeout
		return c.PostLoad()
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
