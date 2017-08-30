package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/efritz/nacelle"
	"github.com/gorilla/mux"

	"github.com/efritz/nacelle/example/api"
)

type Server struct {
	Logger        nacelle.Logger    `service:"logger"`
	SecretService api.SecretService `service:"secret_service"`
	port          int
	server        *http.Server
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

	handlers := &handlerSet{
		logger:        s.Logger,
		secretService: s.SecretService,
	}

	router := mux.NewRouter()
	router.HandleFunc("/post", handlers.post).Methods("POST")
	router.HandleFunc("/load/{id}", handlers.load).Methods("GET")

	s.port = serverConfig.HTTPPort
	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	s.server = &http.Server{Addr: addr, Handler: router}
	return nil
}

func (s *Server) Start() error {
	s.Logger.Info(nil, "HTTP server listening on port %d", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	s.Logger.Info(nil, "Stopping HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.server.Shutdown(ctx)
}
