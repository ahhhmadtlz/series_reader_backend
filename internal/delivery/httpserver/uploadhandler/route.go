package uploadhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	// Protected routes - require authentication
	protectedGroup := e.Group("")
	protectedGroup.Use(middleware.Auth(authService, authConfig))
	protectedGroup.Use(middleware.UserContext())

	// Upload avatar - any authenticated user can upload their own avatar
	protectedGroup.POST("/users/avatar", h.uploadAvatar)

	// Upload series cover - requires manager or admin role
	protectedGroup.POST("/series/:id/cover", h.uploadCover, middleware.RequireManagerOrAdmin())
}