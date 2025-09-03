package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func (a *appHandler) authenticate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &requestPayload)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	auth, err := a.userService.Authentication(ctx, requestPayload.Email, requestPayload.Password)
	if err != nil || !auth {
		errorJSON(w, fmt.Errorf("invalid credentials: %w", err), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
	}

	writeJSON(w, http.StatusOK, payload)
}
