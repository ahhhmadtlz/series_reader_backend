package chapterhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) delete(c echo.Context) error {
	seriesSlug := c.Param("series_slug")
	chapterNumStr:= c.Param("chapter_number")

	chapterNum, err := strconv.ParseFloat(chapterNumStr,64)
	if err !=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chpater number",
		})
	}

	series, err:=h.seriesService.GetByFullSlug(c.Request().Context(),seriesSlug)
	if err !=nil {
		return httpmsgerrorhandler.Error(c,err)
	}

	chapters,err:=h.chapterService.GetBySeriesID(c.Request().Context(),series.ID)
	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	var chapterID uint

	for _,ch:=range chapters {
		if ch.ChapterNumber==chapterNum{
			chapterID=ch.ID
			break
		}
	}
	if chapterID==0 {
		return c.JSON(http.StatusNotFound,httpmsgerrorhandler.ErrorResponse{
			Message: "chapter not found",
		})
	}

	err =h.chapterService.Delete(c.Request().Context(),chapterID)
	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	return  c.NoContent(http.StatusNoContent)
}
