package process

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/efritz/nacelle"
)

type (
	HTTPServer struct {
		Logger          nacelle.Logger            `service:"logger"`
		Container       *nacelle.ServiceContainer `service:"container"`
		configToken     interface{}
		initializer     HTTPServerInitializer
		listener        *net.TCPListener
		server          *http.Server
		once            *sync.Once
		port            int
		certFile        string
		keyFile         string
		shutdownTimeout time.Duration
	}

	HTTPServerInitializer interface {
		Init(nacelle.Config, *http.Server) error
	}

	HTTPServerInitializerFunc func(nacelle.Config, *http.Server) error
)

var ErrBadHTTPConfig = errors.New("HTTP config not registered properly")

func (f HTTPServerInitializerFunc) Init(config nacelle.Config, server *http.Server) error {
	return f(config, server)
}

func NewHTTPServer(initializer HTTPServerInitializer, configs ...HTTPServerConfigFunc) *HTTPServer {
	options := getHTTPOptions(configs)

	return &HTTPServer{
		configToken: options.configToken,
		initializer: initializer,
		once:        &sync.Once{},
	}
}

func (s *HTTPServer) Init(config nacelle.Config) (err error) {
	httpConfig := &HTTPConfig{}
	if err = config.Fetch(s.configToken, httpConfig); err != nil {
		return ErrBadHTTPConfig
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

func (s *HTTPServer) Start() error {
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

func (s *HTTPServer) Stop() (err error) {
	s.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		s.Logger.Info("Shutting down HTTP server")
		err = s.server.Shutdown(ctx)
	})

	return
}
