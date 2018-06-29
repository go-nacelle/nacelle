package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/efritz/nacelle"
)

type (
	Server struct {
		Logger          nacelle.Logger           `service:"logger"`
		Container       nacelle.ServiceContainer `service:"container"`
		configToken     interface{}
		initializer     ServerInitializer
		listener        *net.TCPListener
		server          *http.Server
		once            *sync.Once
		port            int
		certFile        string
		keyFile         string
		shutdownTimeout time.Duration
	}

	ServerInitializer interface {
		Init(nacelle.Config, *http.Server) error
	}

	ServerInitializerFunc func(nacelle.Config, *http.Server) error
)

var ErrBadConfig = fmt.Errorf("HTTP config not registered properly")

func (f ServerInitializerFunc) Init(config nacelle.Config, server *http.Server) error {
	return f(config, server)
}

func NewServer(initializer ServerInitializer, configs ...ConfigFunc) *Server {
	options := getOptions(configs)

	return &Server{
		configToken: options.configToken,
		initializer: initializer,
		once:        &sync.Once{},
	}
}

func (s *Server) Init(config nacelle.Config) (err error) {
	httpConfig := &Config{}
	if err = config.Fetch(s.configToken, httpConfig); err != nil {
		return ErrBadConfig
	}

	s.listener, err = makeListener(httpConfig.HTTPPort)
	if err != nil {
		return err
	}

	s.server = &http.Server{}
	s.port = httpConfig.HTTPPort
	s.certFile = httpConfig.HTTPCertFile
	s.keyFile = httpConfig.HTTPKeyFile
	s.shutdownTimeout = httpConfig.ShutdownTimeout

	if err := s.Container.Inject(s.initializer); err != nil {
		return err
	}

	return s.initializer.Init(config, s.server)
}

func (s *Server) Start() error {
	defer s.listener.Close()
	defer s.server.Close()

	if s.certFile == "" {
		s.Logger.Info("Serving HTTP on port %d", s.port)
		if err := s.server.Serve(s.listener); err != http.ErrServerClosed {
			return err
		}

		s.Logger.Info("No longer serving HTTP on port %d", s.port)
		return nil
	}

	s.Logger.Info("Serving HTTP/TLS on port %d", s.port)
	if err := s.server.ServeTLS(s.listener, s.certFile, s.keyFile); err != http.ErrServerClosed {
		return err
	}

	s.Logger.Info("No longer serving HTTP/TLS on port %d", s.port)
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

func makeListener(port int) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}
