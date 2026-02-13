package chapterhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) create(c echo.Context) error {
 seriesSlug:=c.Param("series_slug")

 series,err:=h.seriesService.GetByFullSlug(c.Request().Context(),seriesSlug)
 if err !=nil{
	return httpmsgerrorhandler.Error(c,err)
 }
 var req param.CreateChapterRequest
 if err:=c.Bind(&req);err!=nil{
	return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
		Message: "invalid request body",
	})
 }
 req.SeriesID=series.ID

 if err := h.chapterValidator.ValidateCreateChapterRequest(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
 }

 response,err:=h.chapterService.Create(c.Request().Context(),req)

 if err !=nil{
	return httpmsgerrorhandler.Error(c,err)
 }

 return c.JSON(http.StatusCreated,response)
}