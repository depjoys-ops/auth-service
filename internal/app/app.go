package app

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/depjoys-ops/auth-service/db"
	"github.com/depjoys-ops/auth-service/internal/config"
	"github.com/depjoys-ops/auth-service/internal/delivery/httpserver"
	"github.com/depjoys-ops/auth-service/internal/repository"
	"github.com/depjoys-ops/auth-service/internal/service"
	"github.com/depjoys-ops/auth-service/pkg/database"
	"github.com/depjoys-ops/auth-service/pkg/logger"
)

func Run(cfg *config.Config, log logger.Logger) {

	if err := db.RunMigrations(cfg.Databases["users"].MigrateUrl); err != nil {
		log.Error("migration failed", err)
		return
	} else {
		log.Info("migrations applied successfully")
	}

	if v, err := db.CheckMigrations(cfg.Databases["users"].MigrateUrl); err != nil {
		log.Error("check migration failed", err)
		return
	} else {
		log.Info("database at migration", "version", v)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := database.NewPostgresPool(ctx, cfg.Databases["users"].URL)
	if err != nil {
		log.Error("can not create DB pool", err)
		return
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
	serverErrChan := httpServer.Start()

	select {
	case <-ctx.Done():
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return
		}

	case err := <-serverErrChan:
		if err != nil {
			return
		}
	}

}
