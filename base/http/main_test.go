package http

//go:generate go-mockgen github.com/efritz/nacelle/config -i Config -o mock_config_test.go -f

import (
	"net"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/config/tag"
	"github.com/efritz/nacelle/service"
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
	config := NewMockConfig()
	config.LoadFunc = func(target interface{}, modifiers ...tag.TagModifier) error {
		c := target.(*Config)
		c.HTTPPort = base.HTTPPort
		c.HTTPCertFile = base.HTTPCertFile
		c.HTTPKeyFile = base.HTTPKeyFile
		c.RawShutdownTimeout = base.RawShutdownTimeout
		return c.PostLoad()
	}

	return config
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
