package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/depjoys-ops/auth-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		pool: pool,
	}
}

var _ domain.UserRepository = (*postgresRepository)(nil)

func (r *postgresRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, first_name, last_name, password, activated)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, created_at, version`

	args := []any{user.Email, user.FirstName, user.LastName, user.Password.Hash, user.Activated}
	err := r.pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%w: %s", domain.ErrEmailAlreadyExists, pgErr.Detail)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {

	query := `SELECT id, email, first_name, last_name, password, activated,	created_at, updated_at
				FROM users
				WHERE id = $1`

	var user domain.User
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password.Hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, first_name, last_name, password, activated,	created_at, updated_at
				FROM users
				WHERE email = $1`

	var user domain.User
	row := r.pool.QueryRow(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password.Hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *postgresRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, first_name = $2, last_name = $3, password = $4, activated = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version`

	args := []any{
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	err := r.pool.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%w: %s", domain.ErrEmailAlreadyExists, pgErr.Detail)
			}
		}

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrRecordNotFound
	}

	return nil
}

func (r *postgresRepository) List(ctx context.Context) ([]*domain.User, error) {
	query := `
			SELECT id, email, first_name, last_name, activated, created_at, updated_at, version
			FROM users`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Activated,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}
