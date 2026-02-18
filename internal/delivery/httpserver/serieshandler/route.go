package serieshandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo, authService auth.Service, authConfig auth.Config) {
	publicGroup := e.Group("/series")
	publicGroup.GET("",h.getList)
  publicGroup.GET("/:identifier", h.get) 


	protectedGroup := e.Group("/series")
	protectedGroup.Use(middleware.Auth(authService,authConfig))
	
	protectedGroup.POST("", h.create,middleware.RequireManagerOrAdmin())
	protectedGroup.PUT("/:id",h.update,middleware.RequireManagerOrAdmin())
	protectedGroup.DELETE("/:id",h.delete,middleware.RequireAdmin())

}