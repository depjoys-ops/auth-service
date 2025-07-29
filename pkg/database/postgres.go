package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	if config.MaxConns == 0 {
		config.MaxConns = 25
	}

	if config.MinConns == 0 {
		config.MinConns = 5
	}

	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = time.Hour
	}

	if config.MaxConnIdleTime == 0 {
		config.MaxConnIdleTime = 30 * time.Minute
	}

	if config.HealthCheckPeriod == 0 {
		config.HealthCheckPeriod = time.Minute
	}

	if config.ConnConfig.ConnectTimeout == 0 {
		config.ConnConfig.ConnectTimeout = 3 * time.Second
	}

	var applicationName = "auth-service"
	var timezone = "UTC"
	if val, ok := config.ConnConfig.RuntimeParams["application_name"]; ok {
		applicationName = val
	}
	if val, ok := config.ConnConfig.RuntimeParams["timezone"]; ok {
		timezone = val
	}
	config.ConnConfig.RuntimeParams = map[string]string{
		"application_name": applicationName,
		"timezone":         timezone,
	}

	return pgxpool.NewWithConfig(ctx, config)
}
