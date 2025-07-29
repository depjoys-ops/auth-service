package main

import (
	"github.com/depjoys-ops/auth-service/internal/app"
	"github.com/depjoys-ops/auth-service/internal/config"
	"github.com/depjoys-ops/auth-service/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger, err := logger.NewSlogLogger(cfg.Logger.FilePath, cfg.Logger.Level)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	app.Run(cfg, logger)
}
