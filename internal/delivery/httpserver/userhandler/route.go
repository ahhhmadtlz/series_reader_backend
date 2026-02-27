package userhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	publicGroup := e.Group("/users")

	// Tight IP limit on auth endpoints — 3 req/sec burst 5.
	// Prevents brute force and credential stuffing.
	authLimiter := middleware.IPRateLimit(3, 5)
	publicGroup.POST("/register", h.register, authLimiter)
	publicGroup.POST("/login", h.login, authLimiter)
	publicGroup.POST("/refresh", h.refresh, authLimiter)

	protectedGroup := e.Group("/users")
	protectedGroup.Use(middleware.Auth(authService, authConfig))
	protectedGroup.Use(middleware.UserContext())
	protectedGroup.Use(middleware.UserRateLimit(30, 50))

	protectedGroup.GET("/profile", h.getProfile)
	protectedGroup.POST("/logout", h.logout)
	protectedGroup.PUT("/profile", h.updateProfile)
	protectedGroup.PUT("/password", h.changePassword)
}