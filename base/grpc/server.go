package grpc

import (
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/efritz/nacelle"
)

type (
	Server struct {
		Logger        nacelle.Logger           `service:"logger"`
		Container     nacelle.ServiceContainer `service:"container"`
		configToken   interface{}
		initializer   ServerInitializer
		listener      *net.TCPListener
		server        *grpc.Server
		once          *sync.Once
		port          int
		serverOptions []grpc.ServerOption
	}

	ServerInitializer interface {
		Init(nacelle.Config, *grpc.Server) error
	}

	ServerInitializerFunc func(nacelle.Config, *grpc.Server) error
)

var ErrBadConfig = fmt.Errorf("gRPC config not registered properly")

func (f ServerInitializerFunc) Init(config nacelle.Config, server *grpc.Server) error {
	return f(config, server)
}

func NewServer(initializer ServerInitializer, configs ...ConfigFunc) *Server {
	options := getOptions(configs)

	return &Server{
		configToken:   options.configToken,
		initializer:   initializer,
		once:          &sync.Once{},
		serverOptions: options.serverOptions,
	}
}

func (s *Server) Init(config nacelle.Config) (err error) {
	grpcConfig := &Config{}
	if err = config.Fetch(s.configToken, grpcConfig); err != nil {
		return ErrBadConfig
	}

	s.listener, err = makeListener(grpcConfig.GRPCPort)
	if err != nil {
		return
	}

	if err := s.Container.Inject(s.initializer); err != nil {
		return err
	}

	s.port = grpcConfig.GRPCPort
	s.server = grpc.NewServer(s.serverOptions...)
	err = s.initializer.Init(config, s.server)
	return
}

func (s *Server) Start() error {
	defer s.listener.Close()

	s.Logger.Info("Serving gRPC on port %d", s.port)

	if err := s.server.Serve(s.listener); err != nil {
		return err
	}

	s.Logger.Info("No longer serving gRPC on port %d", s.port)
	return nil
}

func (s *Server) Stop() error {
	s.once.Do(func() {
		s.Logger.Info("Shutting down gRPC server")
		s.server.GracefulStop()
	})

	return nil
}

func makeListener(port int) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}
