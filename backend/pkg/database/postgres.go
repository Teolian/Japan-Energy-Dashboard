// Package database provides PostgreSQL connection and management.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps sql.DB with utility methods
type DB struct {
	*sql.DB
}

// Config holds database connection configuration
type Config struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConfig returns sensible defaults for Railway PostgreSQL
func DefaultConfig() Config {
	return Config{
		URL:             os.Getenv("DATABASE_URL"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// Connect establishes a connection to PostgreSQL
func Connect(config Config) (*DB, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// RunMigrations executes SQL migrations
func (db *DB) RunMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS energy_data (
			id SERIAL PRIMARY KEY,
			data_type VARCHAR(50) NOT NULL,  -- 'demand', 'jepx', 'reserve'
			area VARCHAR(50),                -- 'tokyo', 'kansai', NULL for system-wide
			date DATE NOT NULL,
			data JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(data_type, area, date)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_energy_data_lookup ON energy_data(data_type, area, date)`,
		`CREATE INDEX IF NOT EXISTS idx_energy_data_date ON energy_data(date DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_energy_data_jsonb ON energy_data USING GIN (data)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// HealthCheck verifies database connectivity
func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
