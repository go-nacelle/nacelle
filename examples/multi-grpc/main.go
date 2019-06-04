package main

//go:generate protoc hello.proto --go_out=plugins=grpc:. -I.

import (
	"github.com/go-nacelle/nacelle"
	basegrpc "github.com/go-nacelle/nacelle/base/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type helloService struct{}
type howdyService struct{}

func (s *helloService) Hello(ctx context.Context, in *empty.Empty) (*HelloResponse, error) {
	return &HelloResponse{Message: "Hello"}, nil
}

func (s *howdyService) Howdy(ctx context.Context, in *empty.Empty) (*HowdyResponse, error) {
	return &HowdyResponse{Message: "Howdy"}, nil
}

func setupServerA(config nacelle.Config, server *grpc.Server) error {
	RegisterHelloServiceServer(server, &helloService{})
	return nil
}

func setupServerB(config nacelle.Config, server *grpc.Server) error {
	RegisterHowdyServiceServer(server, &howdyService{})
	return nil
}

//
//

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	serverA := basegrpc.NewServer(basegrpc.ServerInitializerFunc(setupServerA))
	serverB := basegrpc.NewServer(basegrpc.ServerInitializerFunc(setupServerB), basegrpc.WithTagModifiers(
		nacelle.NewEnvTagPrefixer("b"),
		nacelle.NewDefaultTagSetter("GRPCPort", "6001"),
	))

	processes.RegisterProcess(serverA, nacelle.WithProcessName("grpc-a"))
	processes.RegisterProcess(serverB, nacelle.WithProcessName("grpc-b"))
	return nil
}

//
//

func main() {
	nacelle.NewBootstrapper("multi-grpc-example", setup).BootAndExit()
}
