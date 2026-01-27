package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/QuUteO/video-communication/internal/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func NewClient(ctx context.Context, ps *config.Postgres) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s", ps.User, ps.Password, ps.Host, ps.Port, ps.DB)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	pool, err = pgxpool.Connect(ctx, dsn)
	defer cancel()
	if err == nil {
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to database: %w", err)
}
