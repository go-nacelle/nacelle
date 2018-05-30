package process

import (
	"errors"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/efritz/nacelle"
)

type (
	GRPCServer struct {
		Logger        nacelle.Logger            `service:"logger"`
		Container     *nacelle.ServiceContainer `service:"container"`
		configToken   interface{}
		initializer   GRPCServerInitializer
		listener      *net.TCPListener
		server        *grpc.Server
		once          *sync.Once
		port          int
		serverOptions []grpc.ServerOption
	}

	GRPCServerInitializer interface {
		Init(nacelle.Config, *grpc.Server) error
	}

	GRPCServerInitializerFunc func(nacelle.Config, *grpc.Server) error
)

var ErrBadGRPCConfig = errors.New("gRPC config not registered properly")

func (f GRPCServerInitializerFunc) Init(config nacelle.Config, server *grpc.Server) error {
	return f(config, server)
}

func NewGRPCServer(initializer GRPCServerInitializer, configs ...GRPCServerConfigFunc) *GRPCServer {
	options := getGRPCOptions(configs)

	return &GRPCServer{
		configToken:   options.configToken,
		initializer:   initializer,
		once:          &sync.Once{},
		serverOptions: options.serverOptions,
	}
}

func (s *GRPCServer) Init(config nacelle.Config) (err error) {
	grpcConfig := &GRPCConfig{}
	if err = config.Fetch(s.configToken, grpcConfig); err != nil {
		return ErrBadGRPCConfig
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

func (s *GRPCServer) Start() error {
	defer s.listener.Close()

	s.Logger.Info("Serving gRPC on port %d", s.port)

	if err := s.server.Serve(s.listener); err != nil {
		return err
	}

	s.Logger.Info("No longer serving gRPC on port %d", s.port)
	return nil
}

func (s *GRPCServer) Stop() error {
	s.once.Do(func() {
		s.Logger.Info("Shutting down gRPC server")
		s.server.GracefulStop()
	})

	return nil
}
