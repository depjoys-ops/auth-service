package service

import (
	"context"

	"github.com/depjoys-ops/auth-service/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email string) (*domain.User, error)
	Authentication(ctx context.Context, email, password string) (bool, error)
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

func (u *userService) Authentication(ctx context.Context, email, password string) (bool, error) {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	checked, err := user.Password.CheckPassword(password)
	if err != nil {
		return false, err
	}

	if checked {
		return true, nil
	}

	return false, nil
}
