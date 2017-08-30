package grpc

import (
	"fmt"
	"net"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/example/api"
	"google.golang.org/grpc"
)

type Server struct {
	Logger        nacelle.Logger    `service:"logger"`
	SecretService api.SecretService `service:"secret_service"`
	listener      net.Listener
	grpcServer    *grpc.Server
	port          int
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(config nacelle.Config) error {
	cfg, err := config.Get(ConfigToken)
	if err != nil {
		return err
	}

	serverConfig := cfg.(*Config)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", serverConfig.GRPCPort))
	if err != nil {
		return err
	}

	s.listener = listener
	s.grpcServer = grpc.NewServer()
	RegisterSecretServiceServer(s.grpcServer, s)
	s.port = serverConfig.GRPCPort
	return nil
}

func (s *Server) Start() error {
	defer s.listener.Close()
	s.Logger.Info(nil, "gRPC server listening on port %d", s.port)
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() error {
	s.Logger.Info(nil, "Stopping gRPC server")
	s.grpcServer.GracefulStop()
	return nil
}
