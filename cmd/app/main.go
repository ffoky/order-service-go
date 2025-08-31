package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var (
	PostgresUser     = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresPort     = os.Getenv("POSTGRES_PORT")
	PostgresDB       = os.Getenv("POSTGRES_DB")
)

func newPostgresConnection(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		PostgresUser, PostgresPassword, PostgresPort, PostgresDB,
	)
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool connect %w", err)
	}
	return pool, nil
}

func main() {
	ctx := context.Background()

	pool, err := newPostgresConnection(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

}
