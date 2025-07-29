package httpserver

import (
	"github.com/depjoys-ops/auth-service/internal/service"
	"github.com/depjoys-ops/auth-service/pkg/logger"
)

type appHandler struct {
	userService service.UserService
	logger      logger.Logger
}

func NewAppHandler(service service.UserService, logger logger.Logger) *appHandler {
	return &appHandler{
		userService: service,
		logger:      logger,
	}
}
