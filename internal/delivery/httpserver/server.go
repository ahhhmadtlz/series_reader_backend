package httpserver

import (
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config config.Config
	Router *echo.Echo
}

// New creates a new HTTP server
func New(cfg config.Config) *Server {
	return &Server{
		Router: echo.New(),
		config: cfg,
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