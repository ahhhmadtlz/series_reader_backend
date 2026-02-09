package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/repository/migrator"
	"github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load(config.DefaultOption())
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize GLOBAL logger FIRST (critical!)
	// All other components will use logger.Info(), logger.Error(), etc.
	appLogger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	logger.Info("Application starting",
		"version", "1.0.0",
		"port", cfg.HTTPServer.Port,
	)

	// 3. Initialize database (now logger is ready)
	postgresDB, err := postgres.New(cfg.Postgres)
	if err != nil {
		logger.Fatal("Failed to connect to database",
			"error", err.Error(),
		)
	}

	// 4. Run migrations
	mg := migrator.New(postgresDB.Conn())
	if err := mg.Up(); err != nil {
		logger.Fatal("Failed to run migrations",
			"error", err.Error(),
		)
	}
	logger.Info("Migrations completed successfully")

	server := httpserver.New(*cfg)

	
	logger.Info("Server configured successfully",
		"routes", "health-check only",
	)

	// 6. Start server in a goroutine
	go func() {
		server.Serve()
	}()

	logger.Info("Test the health check at: http://localhost:8080/health-check")

	// 7. Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Received interrupt signal, shutting down gracefully...")

	// 8. Close database connection
	if err := postgresDB.Close(); err != nil {
		logger.Error("Failed to close database connection",
			"error", err.Error(),
		)
	}

	logger.Info("Application stopped")
}