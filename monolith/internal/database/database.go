package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/go-wallet/monolith/internal/logger"
)

const (
	dbMaxRetries  = 5
	dbBaseBackoff = time.Second
)

func ConnectWithRetry(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	var (
		pool *pgxpool.Pool
		err  error
	)

	for attempt := 1; attempt <= dbMaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("connection canceled: %w", ctx.Err())
		default:
		}

		pool, err = pgxpool.New(ctx, connectionString)
		if err != nil {
			if attempt == dbMaxRetries {
				return nil, fmt.Errorf("failed to create pgxpool after %d attempts: %+v", attempt, err)
			}
			time.Sleep(backoffDuration(attempt))
			continue
		}

		pingCtx, cancel := context.WithTimeout(ctx, dbBaseBackoff*2)

		err = pool.Ping(pingCtx)
		cancel()

		if err == nil {
			logger.Log.Info(fmt.Sprintf("Database connected successfully on attempt %d", attempt))
			return pool, nil
		}

		pool.Close()

		if attempt == dbMaxRetries {
			return nil, fmt.Errorf("failed to ping database after %d attempts: %+v", attempt, err)
		}

		logger.Log.Warn(fmt.Sprintf("Database connection attempt %d/%d failed %w", attempt, dbMaxRetries, err))
		time.Sleep(backoffDuration(attempt))
	}
	return nil, fmt.Errorf("failed to connect to database after %d attempts", dbMaxRetries)
}

func backoffDuration(attempt int) time.Duration {
	return dbBaseBackoff * time.Duration(1<<attempt)
}
