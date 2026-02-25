package imageprocessinghandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	protectedGroup := e.Group("")
	protectedGroup.Use(middleware.Auth(authService, authConfig))
	protectedGroup.Use(middleware.UserContext())

	// Series variant reads
	protectedGroup.GET("/series/:seriesID/cover-variants", h.getCoverVariants)
	protectedGroup.GET("/series/:seriesID/banner-variants", h.getBannerVariants)

	// Chapter variant reads
	protectedGroup.GET("/chapters/:chapterID/thumbnail-variants", h.getThumbnailVariants)
	protectedGroup.GET("/chapters/:chapterID/pages/:pageID/variants", h.getPageVariants)
}