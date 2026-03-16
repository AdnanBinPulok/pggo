package connection

import (
	"context"
	"fmt"
)

// Ping verifies the database pool can serve requests.
func (c *DatabaseConnection) Ping() error {
	pool, err := c.Pool()
	if err != nil {
		return err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("db ping failed: %w", err)
	}
	return nil
}
