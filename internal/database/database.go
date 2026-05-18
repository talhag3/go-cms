package database

import (
	"context"
	"fmt"
	"time"

	"github.com/talhag3/go-cms/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a new connection pool to PostgreSQL
//
// ═══════════════════════════════════════════════════════════════════════════
// GO CONCEPT - CONNECTION POOLING:
// ═══════════════════════════════════════════════════════════════════════════
//
// A connection pool maintains a set of ready-to-use database connections.
// Instead of creating a new connection for each request (SLOW),
// we reuse existing connections from the pool (FAST).
//
// PHP COMPARISON:
// - PHP-FPM: Each worker has 1 persistent connection (not a pool)
// - Swoole: Can have a connection pool similar to this
// - Laravel: Uses 1 connection per request (reconnected if needed)
//
// Why pooling in Go?
// - Go handles many requests per process (not like PHP-FPM)
// - Creating DB connections is expensive (~50-100ms each)
// - Pool eliminates this overhead
//
// ┌─────────────────────────────────────────────────────────────────┐
// │                    CONNECTION POOL                              │
// │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                        │
// │  │Conn │ │Conn │ │Conn │ │Conn │ │Conn │  ← Available           │
// │  │  1  │ │  2  │ │  3  │ │  4  │ │  5  │                        │
// │  └──┬──┘ └──┬──┘ └─────┘ └─────┘ └─────┘                        │
// │     │       │                                                   │
// │     ▼       ▼                                                   │
// │  ┌─────┐ ┌─────┐              ← In use by requests              │
// │  │Conn │ │Conn │                                                │
// │  │  6  │ │  7  │                                                │
// │  └─────┘ └─────┘                                                │
// │                                                                 │
// │  When request finishes → connection returns to pool             │
// │  When pool empty + under max → create new connection            │
// │  When pool full + request waits → queue request                 │
// └─────────────────────────────────────────────────────────────────┘
func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	// pgxpool.Config holds all pool configuration
	// We build it from our config struct
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Customize pool settings
	// These are IMPORTANT for production performance
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns

	// MaxConnIdleTime: How long an idle connection stays in pool
	// In PHP: Not applicable (connection dies with request)
	poolConfig.MaxConnIdleTime = time.Duration(cfg.MaxIdleTime) * time.Second

	// MaxConnLifetime: Maximum time ANY connection can exist
	// Prevents using stale connections (database might restart)
	// PRODUCTION: Set to less than your DB's connection timeout
	poolConfig.MaxConnLifetime = time.Duration(cfg.MaxLifetime) * time.Second

	// Health check: pgx will check if connection is alive before giving it
	// This prevents "connection reset" errors
	poolConfig.HealthCheckPeriod = 30 * time.Second

	// Actually create the pool and connect
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// VERIFY CONNECTION WORKS
	// This is like running "SELECT 1" to test connection
	// In PHP: PDO throws exception on bad connect anyway
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// Close gracefully closes the connection pool
// Call this when shutting down the server
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
