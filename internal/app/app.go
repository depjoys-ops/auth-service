package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/depjoys-ops/auth-service/internal/config"
	"github.com/depjoys-ops/auth-service/internal/delivery/httpserver"
	"github.com/depjoys-ops/auth-service/internal/repository"
	"github.com/depjoys-ops/auth-service/internal/service"
	"github.com/depjoys-ops/auth-service/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config, log logger.Logger) {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DBPostgres.Url())
	if err != nil {
		log.Error("can not create pgxpool")
	}
	defer dbPool.Close()

	repository := repository.NewPostgresRepository(dbPool)
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
