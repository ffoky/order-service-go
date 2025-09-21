package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	maxRetries = 5
	baseDelay  = time.Second * 2
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	var pool *pgxpool.Pool
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			lastErr = fmt.Errorf("pgxpool connect: %w", err)
		} else {
			if err := pool.Ping(ctx); err != nil {
				pool.Close()
				lastErr = fmt.Errorf("ping database: %w", err)
			} else {
				return pool, nil
			}
		}

		if attempt < maxRetries {
			delay := time.Duration(attempt) * baseDelay
			logrus.Warnf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...",
				attempt, maxRetries, lastErr, delay)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, lastErr)
}
