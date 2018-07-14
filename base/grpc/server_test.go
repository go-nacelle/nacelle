package grpc

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/base/grpc/internal"
	"github.com/efritz/nacelle/process"
	"github.com/efritz/nacelle/service"
)

type ServerSuite struct{}

func (s *ServerSuite) TestServeAndStop(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		internal.RegisterTestServiceServer(server, &upperService{})

		return nil
	})

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(ConfigToken, &Config{}))
	Expect(err).To(BeNil())

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", getDynamicPort(server.listener)), grpc.WithInsecure())
	Expect(err).To(BeNil())
	defer conn.Close()

	client := internal.NewTestServiceClient(conn)

	resp, err := client.ToUpper(context.Background(), &internal.UpperRequest{Text: "foobar"})
	Expect(err).To(BeNil())
	Expect(resp.GetText()).To(Equal("FOOBAR"))
}

func (s *ServerSuite) TestBadConfig(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		return nil
	})

	err := server.Init(makeConfig(ConfigToken, &emptyConfig{}))
	Expect(err).To(Equal(ErrBadConfig))
}

func (s *ServerSuite) TestBadInjection(t sweet.T) {
	server := NewServer(&badInjectionInitializer{})
	server.Services = makeBadContainer()
	server.Health = process.NewHealth()

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(ConfigToken, &Config{}))
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *ServerSuite) TestInitError(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		return fmt.Errorf("utoh")
	})

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(ConfigToken, &Config{}))
	Expect(err).To(MatchError("utoh"))
}

//
// Helpers

func makeGRPCServer(initializer func(nacelle.Config, *grpc.Server) error) *Server {
	server := NewServer(ServerInitializerFunc(initializer))
	server.Logger = nacelle.NewNilLogger()
	server.Services, _ = service.NewContainer()
	server.Health = process.NewHealth()
	return server
}

//
// Service Impl

type upperService struct{}

func (us *upperService) ToUpper(ctx context.Context, r *internal.UpperRequest) (*internal.UpperResponse, error) {
	return &internal.UpperResponse{Text: strings.ToUpper(r.GetText())}, nil
}

//
// Bad Injection

type badInjectionInitializer struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionInitializer) Init(nacelle.Config, *grpc.Server) error {
	return nil
}
