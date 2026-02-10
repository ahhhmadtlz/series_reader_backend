package serieshandler

import "github.com/labstack/echo/v4"

func (h Handler) SetRoutes(e *echo.Echo) {
	seriesGroup := e.Group("/series")

	seriesGroup.POST("", h.create)
  seriesGroup.GET("/:identifier", h.get) 
	seriesGroup.GET("",h.getList)
	seriesGroup.PUT("/:id",h.update)
	seriesGroup.DELETE("/:id",h.delete)

}