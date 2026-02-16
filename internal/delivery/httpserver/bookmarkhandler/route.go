package bookmarkhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	// All bookmark routes are protected (require authentication)
	bookmarkGroup := e.Group("/bookmarks")
	bookmarkGroup.Use(middleware.Auth(authService, authConfig))
	bookmarkGroup.Use(middleware.UserContext())

	bookmarkGroup.POST("", h.create)
	bookmarkGroup.GET("", h.list)
	bookmarkGroup.DELETE("/:series_id", h.delete)
}