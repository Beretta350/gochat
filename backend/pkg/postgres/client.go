package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
)

// Client wraps the PostgreSQL connection pool
type Client struct {
	Pool *pgxpool.Pool
}

// NewClient creates a new PostgreSQL client (Fx provider)
func NewClient(cfg *config.Config) (*Client, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	// Pool configuration
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	logger.Info("PostgreSQL connected")
	return &Client{Pool: pool}, nil
}

// Close closes the connection pool
func (c *Client) Close() {
	logger.Info("Closing PostgreSQL connection")
	c.Pool.Close()
}
