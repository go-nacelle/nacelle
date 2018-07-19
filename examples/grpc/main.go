package main

//go:generate protoc ping.proto --go_out=plugins=grpc:. -I.

import (
	"github.com/efritz/nacelle"
	basegrpc "github.com/efritz/nacelle/base/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type pingService struct{}

func (s *pingService) Ping(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func setupServer(config nacelle.Config, server *grpc.Server) error {
	RegisterPingServiceServer(server, &pingService{})
	return nil
}

//
//

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterProcess(basegrpc.NewServer(basegrpc.ServerInitializerFunc(setupServer)))
	return nil
}

//
//

func main() {
	nacelle.NewBootstrapper("grpc-example", setup).BootAndExit()
}
