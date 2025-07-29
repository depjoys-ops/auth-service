package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/depjoys-ops/auth-service/internal/config"
	"github.com/depjoys-ops/auth-service/internal/delivery/httpserver"
	"github.com/depjoys-ops/auth-service/internal/repository"
	"github.com/depjoys-ops/auth-service/internal/service"
	"github.com/depjoys-ops/auth-service/pkg/database"
	"github.com/depjoys-ops/auth-service/pkg/logger"
)

func Run(cfg *config.Config, log logger.Logger) {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := database.NewPostgresPool(ctx, cfg.Databases["users"].URL)
	if err != nil {
		log.Error("can not create DB pool")
	}
	defer conn.Close()

	repository := repository.NewPostgresRepository(conn)
	userService := service.NewUserService(repository)
	_ = httpserver.NewAppHandler(userService, log)
	httpRouter := httpserver.NewRouter()

	httpServerOpts := httpserver.ServerOptions{
		Addr:         cfg.HTTPServer.Addr,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ErrorLog:     log.ServerErrorLog(),
	}

	httpServer := httpserver.NewHTTPServer(httpRouter, log, httpServerOpts)
	httpServer.Start()

	<-ctx.Done()
	httpServer.Shutdown(ctx)
}
