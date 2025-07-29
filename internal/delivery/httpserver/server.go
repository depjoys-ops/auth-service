package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/depjoys-ops/auth-service/pkg/logger"
)

type ServerOptions struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ErrorLog     *log.Logger
}

type HTTPServer struct {
	server *http.Server
	logger logger.Logger
}

func NewHTTPServer(handler http.Handler, log logger.Logger, opts ServerOptions) *HTTPServer {
	srv := &http.Server{
		Addr:         opts.Addr,
		Handler:      handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
		ErrorLog:     opts.ErrorLog,
	}

	return &HTTPServer{
		server: srv,
		logger: log,
	}
}

func (s *HTTPServer) Start() {
	s.logger.Info(fmt.Sprintf("HTTP server starting on %s", s.server.Addr))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server failed to start", err.Error())
		}
	}()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.logger.Info("HTTP server shutting down...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP server shutdown failed", err.Error())
		return err
	}
	s.logger.Info("HTTP server stopped")
	return nil
}
