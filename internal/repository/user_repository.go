package repository

import (
	"context"

	"github.com/depjoys-ops/auth-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		db: db,
	}
}

func (p *postgresRepository) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (p *postgresRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return nil, nil
}

func (p *postgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (p *postgresRepository) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (p *postgresRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (p *postgresRepository) List(ctx context.Context) ([]*domain.User, error) {
	return nil, nil
}
