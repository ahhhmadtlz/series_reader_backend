package httpserver

import (
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/adminhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/chapterhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/serieshandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/userhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	seriesservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/service"
	seriesvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"

	chapterservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/service"
	chaptervalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/validator"
	userservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/validator"

	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/bookmarkhandler"
	bookmarkservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/service"
	bookmarkvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/validator"

	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/readinghistoryhandler"
	readinghistoryservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/service"
	readinghistoryvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/validator"
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
}

// New creates a new HTTP server
func New(
	config config.Config,
	authSvc auth.Service,
	seriesSvc seriesservice.Service,
	seriesValidator seriesvalidator.Validator,
	chapterSvc chapterservice.Service,
	chapterValidator chaptervalidator.Validator,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	bookmarkSvc bookmarkservice.Service,       
	bookmarkValidator bookmarkvalidator.Validator, 
	readingHistorySvc readinghistoryservice.Service,
	readingHistoryValidator readinghistoryvalidator.Validator,
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
		adminHandler: adminhandler.New(userSvc),
	}
}

func (s *Server) Serve() {
	// Setup middleware
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.RequestID())

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