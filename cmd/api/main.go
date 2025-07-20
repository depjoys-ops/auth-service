package main

import (
	"github.com/depjoys-ops/auth-service/internal/api"
	"github.com/depjoys-ops/auth-service/internal/config"
)

func main() {
	cfg := config.Load()
	api.Run(cfg)
}
