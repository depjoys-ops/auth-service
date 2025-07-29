package service

import (
	"context"

	"github.com/depjoys-ops/auth-service/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email string) (*domain.User, error)
}

type userService struct {
	repository domain.UserRepository
}

func NewUserService(repository domain.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (u *userService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	return nil, nil
}
