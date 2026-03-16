package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseConnection owns the pgx pool used by all tables.
type DatabaseConnection struct {
	DBURL          string
	MaxConnections int
	Reconnect      bool
	pool           *pgxpool.Pool
}

// Connect initializes the pool if needed and returns it.
func (c *DatabaseConnection) Connect() (*pgxpool.Pool, error) {
	if c.pool != nil {
		return c.pool, nil
	}
	cfg, err := pgxpool.ParseConfig(c.DBURL)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	if c.MaxConnections > 0 {
		cfg.MaxConns = int32(c.MaxConnections)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	c.pool = pool
	return pool, nil
}

// Pool returns a connected pool instance.
func (c *DatabaseConnection) Pool() (*pgxpool.Pool, error) {
	return c.Connect()
}

// Close closes the underlying pool.
func (c *DatabaseConnection) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}
