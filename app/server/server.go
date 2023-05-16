package server

import (
	"context"
	"github.com/shamank/eduTour-backend/app/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler:        handler,
			WriteTimeout:   cfg.HTTP.WriteTimeOut,
			ReadTimeout:    cfg.HTTP.ReadTimeOut,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
