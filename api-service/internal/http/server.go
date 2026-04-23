package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/keshvan/trod-kafka-lab/api-service/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
			Handler: handler,
		},
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on %s", s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
