package uploadhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	protectedGroup := e.Group("")
	protectedGroup.Use(middleware.Auth(authService, authConfig))
	protectedGroup.Use(middleware.UserContext())

	// Avatar — any authenticated user
	protectedGroup.POST("/users/avatar", h.uploadAvatar)

	// Series images — manager or admin
	protectedGroup.POST("/series/:seriesID/cover", h.uploadCover, middleware.RequireManagerOrAdmin())
	protectedGroup.POST("/series/:seriesID/banner", h.uploadBanner, middleware.RequireManagerOrAdmin())

	// Chapter thumbnail — manager or admin
	protectedGroup.POST("/chapters/:chapterID/thumbnail", h.uploadChapterThumbnail, middleware.RequireManagerOrAdmin())
}