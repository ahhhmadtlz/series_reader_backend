package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"

	chapterservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/service"
	chaptervalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/validator"

	seriesservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/service"
	seriesvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"

	userservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/validator"

	bookmarkservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/service"
	bookmarkvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/validator"
	bookmarkrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/bookmark"


	readinghistoryservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/service"
	readinghistoryvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/validator"
	readinghistoryrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/readinghistory"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/repository/migrator"
	"github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres"

	chapterrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/chapter"
	seriesrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/series"
	userrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/user"
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


  // ========================================
	// Phase 2: Setup Services
	// ========================================

authSvc, seriesSvc, seriesValidator, chapterSvc, chapterValidator, userSvc, userValidator, bookmarkSvc, bookmarkValidator, readingHistorySvc, readingHistoryValidator := setupServices(postgresDB, cfg)

  // ========================================
	// Phase 3: HTTP Server Setup
	// ========================================
	server := httpserver.New(
		*cfg,
		authSvc,
		seriesSvc,
		seriesValidator,
		chapterSvc,
		chapterValidator,
		userSvc,
		userValidator,
		bookmarkSvc,
		bookmarkValidator,
		readingHistorySvc,
		readingHistoryValidator,
	)

	logger.Info("HTTP server initialized")

	// 6. Start server in a goroutine
	go func() {
		server.Serve()
	}()

	logger.Info("Server is ready",
		"health_check", "http://localhost:8080/health-check",
		"series_api", "http://localhost:8080/series",
		"auth_api", "http://localhost:8080/users",
	)

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

func setupServices(db *postgres.DB,cfg *config.Config) (
	auth.Service,
	seriesservice.Service,
	seriesvalidator.Validator,
	chapterservice.Service,
	chaptervalidator.Validator,
	userservice.Service,
	uservalidator.Validator,
	bookmarkservice.Service,
	bookmarkvalidator.Validator,
	readinghistoryservice.Service,
	readinghistoryvalidator.Validator,
) {
// ========================================
	// Auth Setup (needed by user service)
	// ========================================
	authService := auth.New(cfg.Auth)
	logger.Info("Auth service initialized")

	// ========================================
	// Series Setup
	// ========================================
	seriesRepository := seriesrepo.New(db.Conn())
	logger.Info("Series repository initialized")

	seriesValidator := seriesvalidator.New()
	logger.Info("Series validator initialized")

	seriesService := seriesservice.New(seriesRepository)
	logger.Info("Series service initialized")

	// ========================================
	// Chapter Setup
	// ========================================
	chapterRepository := chapterrepo.New(db.Conn())
	logger.Info("Chapter repository initialized")

	chapterValidator := chaptervalidator.New()
	logger.Info("Chapter validator initialized")

	chapterService := chapterservice.New(chapterRepository)
	logger.Info("Chapter service initialized")

	// ========================================
	// User Setup
	// ========================================
	userRepository := userrepo.New(db.Conn())
	logger.Info("User repository initialized")

	userValidator := uservalidator.New(userRepository)
	logger.Info("User validator initialized")

	userService := userservice.New(authService, userRepository)
	logger.Info("User service initialized")


	// ========================================
	// Bookmark Setup
	// ========================================
	bookmarkRepository := bookmarkrepo.New(db.Conn())  
	logger.Info("Bookmark repository initialized")

	bookmarkValidator := bookmarkvalidator.New()  
	logger.Info("Bookmark validator initialized")

	bookmarkService := bookmarkservice.New(bookmarkRepository, seriesRepository) 
	logger.Info("Bookmark service initialized")

	// ========================================
	// Reading History Setup
	// ========================================
	readingHistoryRepository := readinghistoryrepo.New(db.Conn())
	logger.Info("Reading history repository initialized")

	readingHistoryValidator := readinghistoryvalidator.New()
	logger.Info("Reading history validator initialized")

	readingHistoryService := readinghistoryservice.New(readingHistoryRepository, chapterRepository, seriesRepository)
	logger.Info("Reading history service initialized")

	return authService, seriesService, seriesValidator, chapterService, chapterValidator, userService, userValidator, bookmarkService, bookmarkValidator, readingHistoryService, readingHistoryValidator
}
