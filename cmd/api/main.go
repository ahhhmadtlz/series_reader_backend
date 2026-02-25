package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/imageprocessor"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage/local"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/worker"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"

	chapterservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/service"
	chaptervalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/validator"

	ipservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/service"

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
	imagevariantrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/imagevariant"
	seriesrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/series"
	userrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/user"

	uploadservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/service"
	uploadvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/validator"
	bannervariantrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/bannervariant"
	chapterthumbnailvariantrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/chapterthumbnailvariant"
	covervariantrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/covervariant"
	uploadrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/upload"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load(config.DefaultOption())
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize logger
	appLogger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	logger.Info("Application starting",
		"version", "1.0.0",
		"port", cfg.HTTPServer.Port,
	)

	// 3. Initialize database
	postgresDB, err := postgres.New(cfg.Postgres)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err.Error())
	}

	// 4. Run migrations
	mg := migrator.New(postgresDB.Conn())
	if err := mg.Up(); err != nil {
		logger.Fatal("Failed to run migrations", "error", err.Error())
	}
	logger.Info("Migrations completed successfully")

	// 5. River migrations
	riverMigrator, err := rivermigrate.New(riverdatabasesql.New(postgresDB.Conn()), nil)
	if err != nil {
		logger.Fatal("Failed to create River migrator", "error", err.Error())
	}
	if _, err := riverMigrator.Migrate(context.Background(), rivermigrate.DirectionUp, nil); err != nil {
		logger.Fatal("Failed to run River migrations", "error", err.Error())
	}
	logger.Info("River tables ready")

	// ========================================
	// Phase 2: Setup Services
	// ========================================
	authSvc, seriesSvc, seriesValidator, chapterSvc, chapterValidator, userSvc, userValidator, bookmarkSvc, bookmarkValidator, readingHistorySvc, readingHistoryValidator, uploadSvc, uploadValidator, imageWorker, jobQueue, ipSvc := setupServices(postgresDB, cfg)

	// ========================================
	// Phase 7: River Worker Setup
	// ========================================
	riverClient, err := worker.NewRiverClient(postgresDB.Conn(), imageWorker)
	if err != nil {
		logger.Fatal("Failed to create River client", "error", err.Error())
	}

	jobQueue.SetClient(riverClient)

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	if err := worker.StartWorker(workerCtx, riverClient); err != nil {
		logger.Fatal("Failed to start River worker", "error", err.Error())
	}
	defer worker.StopWorker(workerCtx, riverClient)

	logger.Info("River worker started")

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
		uploadSvc,
		uploadValidator,
		ipSvc,
	)

	logger.Info("HTTP server initialized")

	go func() {
		server.Serve()
	}()

	logger.Info("Server is ready",
		"health_check", "http://localhost:8080/health-check",
		"series_api", "http://localhost:8080/series",
		"auth_api", "http://localhost:8080/users",
	)

	// 7. Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Received interrupt signal, shutting down gracefully...")

	if err := postgresDB.Close(); err != nil {
		logger.Error("Failed to close database connection", "error", err.Error())
	}

	logger.Info("Application stopped")
}

func setupServices(db *postgres.DB, cfg *config.Config) (
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
	uploadservice.Service,
	uploadvalidator.Validator,
	*worker.ImageProcessingWorker,
	*worker.RiverJobQueue,
	ipservice.Service,
) {
	// ========================================
	// Storage Setup
	// ========================================
	localStorage := local.New(cfg.Upload.BasePath, cfg.Upload.BaseURL)
	logger.Info("Local storage initialized",
		"base_path", cfg.Upload.BasePath,
		"base_url", cfg.Upload.BaseURL,
	)

	// ========================================
	// Auth Setup
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

	// ========================================
	// Image Processing Setup
	// ========================================
	imageVariantRepo := imagevariantrepo.New(db.Conn())
	logger.Info("Image variant repository initialized")

	coverVariantRepo := covervariantrepo.New(db.Conn())
	logger.Info("Cover variant repository initialized")

	bannerVariantRepo := bannervariantrepo.New(db.Conn())
	logger.Info("Banner variant repository initialized")

	chapterThumbnailVariantRepo := chapterthumbnailvariantrepo.New(db.Conn())
	logger.Info("Chapter thumbnail variant repository initialized")

	imageProcessor := imageprocessor.New(localStorage, cfg.Upload.BasePath)
	logger.Info("Image processor initialized")

	jobQueue := worker.NewRiverJobQueue(nil)
	logger.Info("Job queue initialized")

	imageWorker := worker.NewImageProcessingWorker(
			imageProcessor,
			imageVariantRepo,
			coverVariantRepo,
			bannerVariantRepo,
			chapterThumbnailVariantRepo,
			localStorage,
			cfg.Upload.BasePath,
	)
	logger.Info("Image processing worker initialized")

	ipSvc := ipservice.New(imageVariantRepo, coverVariantRepo, bannerVariantRepo, chapterThumbnailVariantRepo)
	logger.Info("Image processing service initialized")

	// ========================================
	// Chapter Setup
	// ========================================
	chapterRepository := chapterrepo.New(db.Conn())
	logger.Info("Chapter repository initialized")

	chapterValidator := chaptervalidator.New(cfg.Upload)
	logger.Info("Chapter validator initialized")

	chapterService := chapterservice.New(chapterRepository, localStorage, jobQueue, imageVariantRepo)

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

	// ========================================
	// Upload Setup
	// ========================================
	uploadRepository := uploadrepo.New(db.Conn())
	logger.Info("Upload repository initialized")

	uploadValidator := uploadvalidator.New(cfg.Upload)
	logger.Info("Upload validator initialized")

	uploadService := uploadservice.New(uploadRepository, userRepository, seriesRepository, chapterRepository, localStorage, jobQueue)
	logger.Info("Upload service initialized")

	seriesService := seriesservice.New(seriesRepository, localStorage, uploadRepository, chapterRepository, imageVariantRepo, coverVariantRepo, bannerVariantRepo, chapterThumbnailVariantRepo)
	logger.Info("Series service initialized")

	return authService, seriesService, seriesValidator, chapterService, chapterValidator, userService, userValidator, bookmarkService, bookmarkValidator, readingHistoryService, readingHistoryValidator, uploadService, uploadValidator, imageWorker, jobQueue, ipSvc
}