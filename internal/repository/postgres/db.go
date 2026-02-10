package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	_ "github.com/lib/pq"
)

type DB struct {
	config config.Postgres
	db     *sql.DB
}
  
func (d *DB) Conn() *sql.DB {
	return d.db
}

func (d *DB) Close() error {
	logger.Info("Closing PostgreSQL database connection")

	if err := d.db.Close(); err != nil {
		logger.Error("Failed to close PostgreSQL database",
			"error", err.Error(),
		)
		return err
	}

	logger.Info("PostgreSQL database connection closed successfully")
	return nil
}

// New creates a new PostgreSQL connection using global logger
func New(cfg config.Postgres) (*DB, error) {
	logger.Info("Connecting to PostgreSQL database",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
		"username", cfg.Username,
	)

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open PostgreSQL database",
			"error", err.Error(),
			"host", cfg.Host,
			"port", cfg.Port,
		)
		return nil, fmt.Errorf("can't open postgres db: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping PostgreSQL database",
			"error", err.Error(),
			"host", cfg.Host,
			"port", cfg.Port,
		)
		return nil, fmt.Errorf("can't connect to postgres db: %w", err)
	}

	// Set connection pool settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	logger.Info("Successfully connected to PostgreSQL database",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
		"max_open_conns", 10,
		"max_idle_conns", 10,
		"conn_max_lifetime", "3m",
	)

	return &DB{
		config: cfg,
		db:     db,
	}, nil
}