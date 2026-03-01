package db

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a pgx connection pool using the provided database URL.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
    cfg, err := pgxpool.ParseConfig(url)
    if err != nil {
        return nil, err
    }
    cfg.MaxConns = 5
    cfg.HealthCheckPeriod = 30 * time.Second

    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return nil, err
    }
    return pool, nil
}
