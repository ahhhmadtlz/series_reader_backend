package chapterhandler

import "github.com/labstack/echo/v4"

func (h Handler) SetRoutes(e *echo.Echo) {
	
	e.GET("/series/:series_slug/chapters", h.getList)
	e.POST("/series/:series_slug/chapters", h.create)
	e.GET("/series/:series_slug/chapter/:chapter_number", h.get)
	e.GET("/series/:series_slug/chapter/:chapter_number/read", h.read)
	e.DELETE("/series/:series_slug/chapter/:chapter_number", h.delete)

	e.POST("/series/:series_slug/chapter/:chapter_number/pages", h.addPages)
	e.GET("/series/:series_slug/chapter/:chapter_number/pages", h.getPages)
	e.DELETE("/series/:series_slug/chapter/:chapter_number/page/:page_number", h.deletePage)
}
