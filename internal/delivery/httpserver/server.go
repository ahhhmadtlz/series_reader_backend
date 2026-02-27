package httpserver

import (
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/adminhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/chapterhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/imageprocessinghandler"
	custommiddleware "github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/serieshandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/uploadhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/userhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"

	seriesvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"

	chaptervalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/validator"
	userservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/validator"

	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/bookmarkhandler"

	bookmarkvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/validator"

	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/readinghistoryhandler"

	readinghistoryvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/validator"

	uploadvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/validator"
)

type Server struct {
	config config.Config
	Router *echo.Echo
	authService auth.Service
	seriesHandler serieshandler.Handler
	chapterHandler chapterhandler.Handler
	userHandler userhandler.Handler
	bookmarkHandler bookmarkhandler.Handler
	readingHistoryHandler readinghistoryhandler.Handler
	adminHandler adminhandler.Handler
	uploadHandler uploadhandler.Handler
	imageProcessingHandler imageprocessinghandler.Handler
}

// New creates a new HTTP server
func New(
   config config.Config,
    authSvc auth.Service,
    seriesSvc serieshandler.SeriesService,
    seriesValidator seriesvalidator.Validator,
    chapterSvc chapterhandler.ChapterService,
    chapterSeriesSvc chapterhandler.SeriesService,
    chapterValidator chaptervalidator.Validator,
    userSvc userhandler.UserService,
    userValidator uservalidator.Validator,
    bookmarkSvc bookmarkhandler.BookmarkService,
    bookmarkValidator bookmarkvalidator.Validator,
    readingHistorySvc readinghistoryhandler.ReadingHistoryService,
    readingHistoryValidator readinghistoryvalidator.Validator,
    uploadSvc uploadhandler.UploadService,
    uploadValidator uploadvalidator.Validator,
    ipSvc imageprocessinghandler.ImageProcessingService,
    adminSvc userservice.Service,
) Server {
	return Server{
		Router: echo.New(),
		config: config,
		authService: authSvc,
		seriesHandler: serieshandler.New(seriesSvc, seriesValidator),
		chapterHandler: chapterhandler.New(chapterSvc, seriesSvc, chapterValidator),
		userHandler: userhandler.New(userSvc,userValidator),
		bookmarkHandler: bookmarkhandler.New(bookmarkSvc, bookmarkValidator),
		readingHistoryHandler: readinghistoryhandler.New(readingHistorySvc, readingHistoryValidator),
		adminHandler: adminhandler.New(adminSvc),
		uploadHandler: uploadhandler.New(uploadSvc,uploadValidator),
		imageProcessingHandler: imageprocessinghandler.New(ipSvc),
	}
}

func (s *Server) Serve() {
	// Body limit — reject oversized requests before any handler reads the body.
	// Set to the largest expected upload (max_page_size_mb=15) plus headroom.
	s.Router.Use(middleware.BodyLimit(fmt.Sprintf("%dMB", s.config.HTTPServer.BodyLimitMB)))

	// CORS — must be before any route handler so preflight OPTIONS requests
	// are handled before auth middleware rejects them.
	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.config.CORS.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge: 86400, // cache preflight for 24h
	}))
	
	// Setup middleware
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.RequestID())


	// Global IP rate limit — covers all routes including public ones.
	// 10 req/sec per IP, burst of 20. Raised cost for scrapers and bots.
	// See ratelimit.go for the NOTE on IPExtractor if behind a proxy.
	s.Router.Use(custommiddleware.IPRateLimit(10, 20))

	// Request logger middleware using global logger
	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogError:         true,
		LogProtocol:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info("request",
					"request_id", v.RequestID,
					"host", v.Host,
					"content_length", v.ContentLength,
					"protocol", v.Protocol,
					"method", v.Method,
					"latency", v.Latency,
					"remote_ip", v.RemoteIP,
					"response_size", v.ResponseSize,
					"uri", v.URI,
					"status", v.Status,
				)
			} else {
				logger.Error("request",
					"request_id", v.RequestID,
					"host", v.Host,
					"content_length", v.ContentLength,
					"protocol", v.Protocol,
					"method", v.Method,
					"latency", v.Latency,
					"error", v.Error.Error(),
					"remote_ip", v.RemoteIP,
					"response_size", v.ResponseSize,
					"uri", v.URI,
					"status", v.Status,
				)
			}
			return nil
		},
	}))

	// Health check endpoint
	s.Router.GET("/health-check", s.healthCheck)

	//register all routes
	s.seriesHandler.SetRoutes(s.Router,s.authService,s.config.Auth)
	s.chapterHandler.SetRoutes(s.Router,s.authService,s.config.Auth)
  s.userHandler.SetRoutes(s.Router, s.authService, s.config.Auth)
	s.bookmarkHandler.SetRoutes(s.Router, s.authService, s.config.Auth)
	s.readingHistoryHandler.SetRoutes(s.Router, s.authService, s.config.Auth)
	s.adminHandler.SetRoutes(s.Router, s.authService, s.config.Auth)
	s.uploadHandler.SetRoutes(s.Router,s.authService,s.config.Auth)
	s.imageProcessingHandler.SetRoutes(s.Router, s.authService, s.config.Auth)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	
	logger.Info("Starting Echo server",
		"address", address,
	)

	if err := s.Router.Start(address); err != nil {
		logger.Error("Failed to start server",
			"error", err.Error(),
		)
	}
}