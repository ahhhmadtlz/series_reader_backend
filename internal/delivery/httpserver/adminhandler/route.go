package adminhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {

	// All admin routes require auth + admin role
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.Auth(authService, authConfig))
	adminGroup.Use(middleware.UserContext())
	adminGroup.Use(middleware.RequireAdmin())

	// User management
	adminGroup.GET("/users/:user_id", h.getUserPermissions)
	adminGroup.PUT("/users/:user_id/role", h.changeUserRole)
	adminGroup.POST("/users/:user_id/permissions", h.grantPermission)
	adminGroup.DELETE("/users/:user_id/permissions", h.revokePermission)
}