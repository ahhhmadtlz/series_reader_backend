package chapterhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	
	// Public routes
	publicGroup := e.Group("/series")
	publicGroup.GET("/:series_slug/chapters", h.getList)
	publicGroup.GET("/:series_slug/chapter/:chapter_number", h.get)
	publicGroup.GET("/:series_slug/chapter/:chapter_number/read", h.read)


 // Protected routes
	protectedGroup := e.Group("/chapters")
	protectedGroup.Use(middleware.Auth(authService, authConfig))
	protectedGroup.Use(middleware.UserContext())



	// Page read
	protectedGroup.GET("/:id/pages", h.getPages)

	// Page write — manager or admin
	protectedGroup.POST("/:id/pages",
		h.uploadPage,
		middleware.RequireManagerOrAdmin(),
	)
	protectedGroup.POST("/:id/pages/bulk",
		h.bulkUploadPages,
		middleware.RequireManagerOrAdmin(),
	)
	protectedGroup.PATCH("/:id/pages/reorder",
		h.reorderPages,
		middleware.RequireManagerOrAdmin(),
	)
	protectedGroup.DELETE("/:id/pages/:page_number",
		h.deletePage,
		middleware.RequireManagerOrAdmin(),
	)


	// Chapter write
	protectedGroup2 := e.Group("/series")
	protectedGroup2.Use(middleware.Auth(authService, authConfig))
	protectedGroup2.Use(middleware.UserContext())

	protectedGroup2.POST("/:series_slug/chapters",
		h.create,
		middleware.RequireManagerOrAdmin(),
	)
	protectedGroup2.DELETE("/:series_slug/chapter/:chapter_number",
		h.delete,
		middleware.RequireAdmin(),
	)
}