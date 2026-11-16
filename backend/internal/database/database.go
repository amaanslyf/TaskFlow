// Package database manages the PostgreSQL connection pool.
//
// A "connection pool" keeps several database connections open
// and ready to use, instead of opening a brand-new connection
// for every single request.  This is much faster.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a connection pool to PostgreSQL.
// The connStr looks something like:
//
//	postgres://user:password@localhost:5432/taskflow?sslmode=disable
func Connect(connStr string) (*pgxpool.Pool, error) {
	// Parse the connection string into a configuration object.
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	// Pool size limits — good defaults for a small app.
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	// Give PostgreSQL 10 seconds to respond during startup.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	// Ping the database to make sure the connection actually works.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return pool, nil
}

