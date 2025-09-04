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
	log    logger.Logger
}

func NewHTTPServer(h *appHandler, log logger.Logger, opts ServerOptions) *HTTPServer {
	srv := &http.Server{
		Addr:         opts.Addr,
		Handler:      newRouter(h),
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
		ErrorLog:     opts.ErrorLog,
	}

	return &HTTPServer{
		server: srv,
		log:    log,
	}
}

func (s *HTTPServer) Start() <-chan error {
	errChan := make(chan error, 1)
	s.log.Info(fmt.Sprintf("HTTP server starting on %s", s.server.Addr))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("HTTP server failed to start", err.Error())
			errChan <- err
		}
	}()

	return errChan
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.log.Info("HTTP server shutting down...")

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("HTTP server shutdown failed", err.Error())
		return err
	}

	s.log.Info("HTTP server stopped")
	return nil
}
