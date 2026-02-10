package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/serieshandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/repository/migrator"
	"github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres"
	seriesrepo "github.com/ahhhmadtlz/series_reader_backend/internal/repository/postgres/series"
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


	// Initialize repositories
	seriesRepository := seriesrepo.New(postgresDB.Conn())
	logger.Info("Series repository initialized")

	// Initialize validators
	seriesValidator := validator.New()
	logger.Info("Series validator initialized")

	// Initialize services
	seriesService := service.New(seriesRepository)
	logger.Info("Series service initialized")

	// Initialize handlers
	seriesHandler := serieshandler.New(seriesService, seriesValidator)
	logger.Info("Series handler initialized")

	// Initialize HTTP server
	server := httpserver.New(*cfg)

	// Setup routes
	seriesHandler.SetRoutes(server.Router)
	logger.Info("Series routes registered",
			"routes", []string{
				"POST /series",
				"GET /series/:id",
				"GET /series/slug/:slug",
				"GET /series",
				"PUT /series/:id",
				"DELETE /series/:id",
			},
		)

	
	

	// 6. Start server in a goroutine
	go func() {
		server.Serve()
	}()

	logger.Info("Server is ready",
		"health_check", "http://localhost:8080/health-check",
		"series_api", "http://localhost:8080/series",
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


// func setupServices(
// 	cfg config.Config,
// 	mysqlDB *mysql.MySQLDB,
// ) (
// 	auth.Service,
// 	userservice.Service,
// 	uservalidator.Validator,
// 	categoryservice.Service,
// 	categoryvalidator.Validator,
// 	transactionservice.Service,
// 	transactionvalidator.Validator,
// ) {
// 	// Auth service
// 	authSvc := auth.New(cfg.Auth)

// 	userRepo := userrepository.New(mysqlDB)
// 	userValidator := uservalidator.New(userRepo)
// 	userSvc := userservice.New(authSvc, userRepo)

// 	categoryRepo := categoryrepository.New(mysqlDB)
// 	categoryValidator := categoryvalidator.New(categoryRepo)
// 	categorySvc := categoryservice.New(categoryRepo)

// 	transactionRepo := transactionrepository.New(mysqlDB)
// 	transactionValidator := transactionvalidator.New(transactionRepo, categoryRepo)
// 	transactionSvc := transactionservice.New(transactionRepo)

// 	return authSvc, userSvc, userValidator, categorySvc, categoryValidator, transactionSvc, transactionValidator
// }