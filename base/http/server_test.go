package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/service"
)

type ServerSuite struct{}

func (s *ServerSuite) TestServeAndStop(t sweet.T) {
	server := makeHTTPServer(func(config nacelle.Config, server *http.Server) error {
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users/foo" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("bar"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
		})

		return nil
	})

	os.Setenv("HTTP_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(ConfigToken, &Config{}))
	Expect(err).To(BeNil())

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	url := fmt.Sprintf("http://localhost:%d/users/foo", getDynamicPort(server.listener))

	req, err := http.NewRequest("GET", url, nil)
	Expect(err).To(BeNil())

	resp, err := http.DefaultClient.Do(req)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	Expect(data).To(Equal([]byte("bar")))
}

func (s *ServerSuite) TestBadConfig(t sweet.T) {
	server := makeHTTPServer(func(config nacelle.Config, server *http.Server) error {
		return nil
	})

	server.Logger = nacelle.NewNilLogger()
	Expect(server.Init(makeConfig(ConfigToken, &emptyConfig{}))).To(Equal(ErrBadConfig))
}

func (s *ServerSuite) TestBadInjection(t sweet.T) {
	server := NewServer(&badInjectionHTTPInitializer{})
	server.Services = makeBadContainer()

	os.Setenv("HTTP_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(ConfigToken, &Config{}))
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *ServerSuite) TestInitError(t sweet.T) {
	server := makeHTTPServer(func(config nacelle.Config, server *http.Server) error {
		return fmt.Errorf("utoh")
	})

	os.Setenv("HTTP_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(ConfigToken, &Config{}))
	Expect(err).To(MatchError("utoh"))
}

//
// Helpers

func makeHTTPServer(initializer func(nacelle.Config, *http.Server) error) *Server {
	server := NewServer(ServerInitializerFunc(initializer))
	server.Logger = nacelle.NewNilLogger()
	server.Services, _ = service.NewContainer()
	return server
}

//
// Bad Injection

type badInjectionHTTPInitializer struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionHTTPInitializer) Init(nacelle.Config, *http.Server) error {
	return nil
}
