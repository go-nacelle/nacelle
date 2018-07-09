package grpc

import (
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/efritz/nacelle"
	uuid "github.com/satori/go.uuid"
)

type (
	Server struct {
		Logger        nacelle.Logger           `service:"logger"`
		Services      nacelle.ServiceContainer `service:"container"`
		Health        nacelle.Health           `service:"health"`
		configToken   interface{}
		initializer   ServerInitializer
		listener      *net.TCPListener
		server        *grpc.Server
		once          *sync.Once
		stopped       chan struct{}
		port          int
		serverOptions []grpc.ServerOption
		healthToken   healthToken
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
		stopped:       make(chan struct{}),
		serverOptions: options.serverOptions,
		healthToken:   healthToken(uuid.Must(uuid.NewV4()).String()),
	}
}

func (s *Server) Init(config nacelle.Config) (err error) {
	if err := s.Health.AddReason(s.healthToken); err != nil {
		return err
	}

	grpcConfig := &Config{}
	if err = config.Fetch(s.configToken, grpcConfig); err != nil {
		return ErrBadConfig
	}

	s.listener, err = makeListener(grpcConfig.GRPCPort)
	if err != nil {
		return
	}

	if err := s.Services.Inject(s.initializer); err != nil {
		return err
	}

	s.port = grpcConfig.GRPCPort
	s.server = grpc.NewServer(s.serverOptions...)
	err = s.initializer.Init(config, s.server)
	return
}

func (s *Server) Start() error {
	defer s.listener.Close()

	if err := s.Health.RemoveReason(s.healthToken); err != nil {
		return err
	}

	s.Logger.Info("Serving gRPC on port %d", s.port)

	if err := s.server.Serve(s.listener); err != nil {
		select {
		case <-s.stopped:
		default:
			return err
		}
	}

	s.Logger.Info("No longer serving gRPC on port %d", s.port)
	return nil
}

func (s *Server) Stop() error {
	s.once.Do(func() {
		s.Logger.Info("Shutting down gRPC server")
		close(s.stopped)
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
