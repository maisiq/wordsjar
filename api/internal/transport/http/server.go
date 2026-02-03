package http

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/logger"
)

type Server struct {
	cfg     config.Config
	server  *http.Server
	handler http.Handler
	log     logger.Logger
}

func NewServer(cfg config.Config, log logger.Logger, handler http.Handler) *Server {
	srv := http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: handler,
	}

	return &Server{
		cfg:     cfg,
		server:  &srv,
		handler: handler,
		log:     log,
	}
}

func (s *Server) Run(ctx context.Context) {
	s.log.Infow("Starting server",
		"addr", s.server.Addr,
	)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Errorw("failed to serve server", "error", err.Error())
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
