package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/config"
)

type (
	Server struct {
		Logger          nacelle.Logger           `service:"logger"`
		Services        nacelle.ServiceContainer `service:"container"`
		Health          nacelle.Health           `service:"health"`
		tagModifiers    []config.TagModifier
		initializer     ServerInitializer
		listener        *net.TCPListener
		server          *http.Server
		once            *sync.Once
		host            string
		port            int
		certFile        string
		keyFile         string
		shutdownTimeout time.Duration
		healthToken     healthToken
	}

	ServerInitializer interface {
		Init(nacelle.Config, *http.Server) error
	}

	ServerInitializerFunc func(nacelle.Config, *http.Server) error
)

func (f ServerInitializerFunc) Init(config nacelle.Config, server *http.Server) error {
	return f(config, server)
}

func NewServer(initializer ServerInitializer, configs ...ConfigFunc) *Server {
	options := getOptions(configs)

	return &Server{
		tagModifiers: options.tagModifiers,
		initializer:  initializer,
		once:         &sync.Once{},
		healthToken:  healthToken(uuid.New().String()),
	}
}

func (s *Server) Init(config nacelle.Config) (err error) {
	if err := s.Health.AddReason(s.healthToken); err != nil {
		return err
	}

	httpConfig := &Config{}
	if err = config.Load(httpConfig, s.tagModifiers...); err != nil {
		return err
	}

	s.listener, err = makeListener(httpConfig.HTTPHost, httpConfig.HTTPPort)
	if err != nil {
		return err
	}

	s.server = &http.Server{}
	s.host = httpConfig.HTTPHost
	s.port = httpConfig.HTTPPort
	s.certFile = httpConfig.HTTPCertFile
	s.keyFile = httpConfig.HTTPKeyFile
	s.shutdownTimeout = httpConfig.ShutdownTimeout

	if err := s.Services.Inject(s.initializer); err != nil {
		return err
	}

	return s.initializer.Init(config, s.server)
}

func (s *Server) Start() error {
	defer s.listener.Close()
	defer s.server.Close()

	if err := s.Health.RemoveReason(s.healthToken); err != nil {
		return err
	}

	if s.certFile != "" {
		return s.serveTLS()
	}

	return s.serve()
}

func (s *Server) serve() error {
	s.Logger.Info("Serving HTTP on %s:%d", s.host, s.port)
	if err := s.server.Serve(s.listener); err != http.ErrServerClosed {
		return err
	}

	s.Logger.Info("No longer serving HTTP on %s:%d", s.host, s.port)
	return nil
}

func (s *Server) serveTLS() error {
	s.Logger.Info("Serving HTTP/TLS on %s:%d", s.host, s.port)
	if err := s.server.ServeTLS(s.listener, s.certFile, s.keyFile); err != http.ErrServerClosed {
		return err
	}

	s.Logger.Info("No longer serving HTTP/TLS on %s:%d", s.host, s.port)
	return nil
}

func (s *Server) Stop() (err error) {
	s.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		s.Logger.Info("Shutting down HTTP server")
		err = s.server.Shutdown(ctx)
	})

	return
}

func makeListener(host string, port int) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}
