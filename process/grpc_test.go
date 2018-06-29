package process

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process/internal"
)

type GRPCSuite struct{}

func (s *GRPCSuite) TestServeAndStop(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		internal.RegisterTestServiceServer(server, &upperService{})

		return nil
	})

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(GRPCConfigToken, &GRPCConfig{}))
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

func (s *GRPCSuite) TestBadConfig(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		return nil
	})

	err := server.Init(makeConfig(GRPCConfigToken, &emptyConfig{}))
	Expect(err).To(Equal(ErrBadGRPCConfig))
}

func (s *GRPCSuite) TestBadInjection(t sweet.T) {
	server := NewGRPCServer(&badInjectionGRPCInitializer{})
	server.Container = makeBadContainer()

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(GRPCConfigToken, &GRPCConfig{}))
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *GRPCSuite) TestInitError(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		return fmt.Errorf("utoh")
	})

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(GRPCConfigToken, &GRPCConfig{}))
	Expect(err).To(MatchError("utoh"))
}

//
// Helpers

func makeGRPCServer(initializer func(nacelle.Config, *grpc.Server) error) *GRPCServer {
	server := NewGRPCServer(GRPCServerInitializerFunc(initializer))
	server.Logger = nacelle.NewNilLogger()
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

type badInjectionGRPCInitializer struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionGRPCInitializer) Init(nacelle.Config, *grpc.Server) error {
	return nil
}
