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
		Container     *nacelle.ServiceContainer `service:"container"`
		initializer   GRPCServerInitializer
		listener      *net.TCPListener
		server        *grpc.Server
		once          *sync.Once
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

func NewGRPCServer(initializer GRPCServerInitializer, serverOptions ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		initializer:   initializer,
		once:          &sync.Once{},
		serverOptions: serverOptions,
	}
}

func (s *GRPCServer) Init(config nacelle.Config) (err error) {
	rawConfig, err := config.Get(GRPCConfigToken)
	if err != nil {
		return err
	}

	grpcConfig, ok := rawConfig.(*GRPCConfig)
	if !ok {
		return ErrBadGRPCConfig
	}

	s.listener, err = makeListener(grpcConfig.GRPCPort)
	if err != nil {
		return
	}

	if err := s.Container.Inject(s.initializer); err != nil {
		return err
	}

	s.server = grpc.NewServer(s.serverOptions...)
	err = s.initializer.Init(config, s.server)
	return
}

func (s *GRPCServer) Start() error {
	defer s.listener.Close()

	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() error {
	s.once.Do(func() { s.server.GracefulStop() })
	return nil
}
