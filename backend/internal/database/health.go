package database

import (
	"context"
	"fmt"
	"time"
)

const healthCheckTimeout = 5 * time.Second

func Health() error {
	db := Get()

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to access postgres sql db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres health check failed: %w", err)
	}

	return nil
}
