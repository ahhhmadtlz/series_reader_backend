package readinghistoryhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	historyGroup := e.Group("/history")
	historyGroup.Use(middleware.Auth(authService, authConfig))
	historyGroup.Use(middleware.UserContext())

	// POST /history - Mark chapter as read
	historyGroup.POST("", h.markAsRead)

	// DELETE /history/:chapter_id - Unmark chapter as read
	historyGroup.DELETE("/:chapter_id", h.unmarkAsRead)

	// GET /history - Get user's reading history (with pagination)
	historyGroup.GET("", h.getUserHistory)

	// GET /series/:slug/progress - Get reading progress for a series
	seriesGroup := e.Group("/series")
	seriesGroup.Use(middleware.Auth(authService, authConfig))
	seriesGroup.Use(middleware.UserContext())
	
	seriesGroup.GET("/:slug/progress", h.getSeriesProgress)
}
