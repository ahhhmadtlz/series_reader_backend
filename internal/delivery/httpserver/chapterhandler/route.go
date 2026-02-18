package chapterhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	
	publicGroup := e.Group("/series")
	publicGroup.GET("/:series_slug/chapters", h.getList)
	publicGroup.GET("/:series_slug/chapter/:chapter_number", h.get)
	publicGroup.GET("/:series_slug/chapter/:chapter_number/read", h.read)
	publicGroup.GET("/:series_slug/chapter/:chapter_number/pages", h.getPages)


	protectedGroup:=e.Group("/series")
	protectedGroup.Use(middleware.Auth(authService, authConfig))
	protectedGroup.Use(middleware.UserContext())


	protectedGroup.POST("/series/:series_slug/chapters",
	 h.create,
	 middleware.RequireManagerOrAdmin(),
	)

	protectedGroup.DELETE("/series/:series_slug/chapter/:chapter_number",
	 h.delete,
	 middleware.RequireAdmin(),
	)
	
	protectedGroup.POST("/series/:series_slug/chapter/:chapter_number/pages",
	 h.addPages,
	 middleware.RequireManagerOrAdmin(),
	)

	e.DELETE("/series/:series_slug/chapter/:chapter_number/page/:page_number",
	 h.deletePage,
	 middleware.RequireAdmin(),
	)
}
