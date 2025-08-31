package order

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"WBTECH_L0/internal/repository/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool Poolx
	sb   squirrel.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: Poolx{pool},
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type Sqlizer interface {
	ToSql() (sql string, args []interface{}, err error)
}

type Poolx struct {
	*pgxpool.Pool
}

func (p *Poolx) getQuerier(ctx context.Context) Querier {
	if tx := postgres.GetTx(ctx); tx != nil {
		return tx
	}
	return p.Pool
}

func (p *Poolx) Getx(ctx context.Context, dest interface{}, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("postgres: to sql: %w", err)
	}

	querier := p.getQuerier(ctx)
	return pgxscan.Get(ctx, querier, dest, query, args...)
}

func (p *Poolx) Selectx(ctx context.Context, dest interface{}, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("postgres: to sql: %w", err)
	}

	querier := p.getQuerier(ctx)
	return pgxscan.Select(ctx, querier, dest, query, args...)
}

func (p *Poolx) Execx(ctx context.Context, sqlizer Sqlizer) (pgconn.CommandTag, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("postgres: to sql: %w", err)
	}

	querier := p.getQuerier(ctx)
	return querier.Exec(ctx, query, args...)
}
