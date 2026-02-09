package migrator

import (
	"database/sql"
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	defaultMigrationTable = "schema_migrations"
	defaultMigrationsDir  = "./internal/repository/postgres/migrations"
)

type Migrator struct {
	db             *sql.DB
	dialect        string
	migrations     *migrate.FileMigrationSource
	migrationTable string
}

type Config struct {
	MigrationsDir  string
	MigrationTable string
}

// New creates a new Migrator instance with an existing DB connection
func New(db *sql.DB) Migrator {
	return NewWithConfig(db, Config{})
}

// NewWithConfig creates a new Migrator instance with custom configuration
func NewWithConfig(db *sql.DB, cfg Config) Migrator {
	migrationsDir := cfg.MigrationsDir
	if migrationsDir == "" {
		migrationsDir = defaultMigrationsDir
	}

	migrationTable := cfg.MigrationTable
	if migrationTable == "" {
		migrationTable = defaultMigrationTable
	}

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	return Migrator{
		db:             db,
		dialect:        "postgres",
		migrations:     migrations,
		migrationTable: migrationTable,
	}
}

// setMigrationTable configures the migration table name
func (m Migrator) setMigrationTable() {
	migrate.SetTable(m.migrationTable)
}

// Up applies all pending migrations
func (m Migrator) Up() error {
	m.setMigrationTable()

	n, err := migrate.Exec(m.db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		logger.Error("Failed to apply migrations",
			"error", err.Error(),
		)
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if n == 0 {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Successfully applied migrations",
			"count", n,
		)
	}

	return nil
}

// UpWithLimit applies up to limit pending migrations
func (m Migrator) UpWithLimit(limit int) error {
	m.setMigrationTable()

	n, err := migrate.ExecMax(m.db, m.dialect, m.migrations, migrate.Up, limit)
	if err != nil {
		logger.Error("Failed to apply migrations",
			"error", err.Error(),
			"limit", limit,
		)
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if n == 0 {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Successfully applied migrations",
			"count", n,
		)
	}

	return nil
}

// Down rolls back the last migration
func (m Migrator) Down() error {
	return m.DownWithLimit(1)
}

// DownWithLimit rolls back up to limit migrations
func (m Migrator) DownWithLimit(limit int) error {
	m.setMigrationTable()

	n, err := migrate.ExecMax(m.db, m.dialect, m.migrations, migrate.Down, limit)
	if err != nil {
		logger.Error("Failed to rollback migrations",
			"error", err.Error(),
			"limit", limit,
		)
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	if n == 0 {
		logger.Info("No migrations to rollback")
	} else {
		logger.Info("Successfully rolled back migrations",
			"count", n,
		)
	}

	return nil
}

// Status returns the current migration status
func (m Migrator) Status() error {
	m.setMigrationTable()

	records, err := migrate.GetMigrationRecords(m.db, m.dialect)
	if err != nil {
		logger.Error("Failed to get migration records",
			"error", err.Error(),
		)
		return fmt.Errorf("failed to get migration records: %w", err)
	}

	migrations, err := m.migrations.FindMigrations()
	if err != nil {
		logger.Error("Failed to find migrations",
			"error", err.Error(),
		)
		return fmt.Errorf("failed to find migrations: %w", err)
	}

	// Create a map of applied migrations
	appliedMap := make(map[string]*migrate.MigrationRecord)
	for _, record := range records {
		appliedMap[record.Id] = record
	}

	logger.Info("\nMigration Status:")
	logger.Info("==================")

	if len(migrations) == 0 {
		logger.Info("No migrations found")
		return nil
	}

	pendingCount := 0
	appliedCount := 0

	for _, migration := range migrations {
		if record, applied := appliedMap[migration.Id]; applied {
			logger.Info("Migration applied",
				"id", migration.Id,
				"applied_at", record.AppliedAt.Format("2006-01-02 15:04:05"),
			)
			appliedCount++
		} else {
			logger.Info("Migration pending",
				"id", migration.Id,
			)
			pendingCount++
		}
	}

	logger.Info("Migration summary",
		"total", len(migrations),
		"applied", appliedCount,
		"pending", pendingCount,
	)

	return nil
}

// Redo rolls back and re-applies the last migration
func (m Migrator) Redo() error {
	logger.Info("Rolling back last migration...")
	if err := m.Down(); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	logger.Info("Re-applying last migration...")
	if err := m.UpWithLimit(1); err != nil {
		return fmt.Errorf("failed to re-apply: %w", err)
	}

	logger.Info("Redo completed successfully")
	return nil
}

// Reset rolls back all migrations
func (m Migrator) Reset() error {
	logger.Info("Rolling back all migrations...")

	m.setMigrationTable()

	// Rollback all by passing 0
	n, err := migrate.Exec(m.db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		logger.Error("Failed to reset migrations",
			"error", err.Error(),
		)
		return fmt.Errorf("failed to reset: %w", err)
	}

	logger.Info("Reset completed successfully",
		"rolled_back", n,
	)
	return nil
}

// Fresh rolls back all migrations and re-applies them
func (m Migrator) Fresh() error {
	if err := m.Reset(); err != nil {
		return err
	}

	logger.Info("Re-applying all migrations...")
	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to re-apply migrations: %w", err)
	}

	logger.Info("Fresh completed successfully")
	return nil
}