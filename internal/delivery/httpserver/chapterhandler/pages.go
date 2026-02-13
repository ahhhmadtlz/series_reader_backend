package chapterhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) addPages(c echo.Context)error {
	seriesSlug:=c.Param("series_slug")
	chapterNumStr:=c.Param("chapter_number")

	chapterNum,err:=strconv.ParseFloat(chapterNumStr,64)
	if err !=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter number",
		})
	}

	series,err:=h.seriesService.GetByFullSlug(c.Request().Context(),seriesSlug)
	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}


	chapters, err:=h.chapterService.GetBySeriesID(c.Request().Context(),series.ID)
	if err!=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	var chapterID uint
  for _,ch:=range chapters {
		if ch.ChapterNumber==chapterNum {
			chapterID =ch.ID
			break
		}
	}

	if chapterID ==0 {
		return c.JSON(http.StatusNotFound,httpmsgerrorhandler.ErrorResponse{
			Message: "chapter not found",
		})
	}

	var req param.AddChapterPagesRequest

	if err:=c.Bind(&req);err!=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}
	req.ChapterID=chapterID

	if err := h.chapterValidator.ValidateAddChapterPagesRequest(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
 }

 if err:=h.chapterService.AddPages(c.Request().Context(),req); err !=nil{
	return httpmsgerrorhandler.Error(c,err)
 }

 return c.JSON(http.StatusCreated,echo.Map{
	"message":"page added succesfully",
 })
}

func (h Handler) getPages(c echo.Context) error {
	seriesSlug := c.Param("series_slug")
	chapterNumStr := c.Param("chapter_number")

	// Parse chapter number
	chapterNum, err := strconv.ParseFloat(chapterNumStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter number",
		})
	}

	// Get series by full_slug
	series, err := h.seriesService.GetByFullSlug(c.Request().Context(), seriesSlug)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Get all chapters for series to find the chapter ID
	chapters, err := h.chapterService.GetBySeriesID(c.Request().Context(), series.ID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Find chapter by number
	var chapterID uint
	for _, ch := range chapters {
		if ch.ChapterNumber == chapterNum {
			chapterID = ch.ID
			break
		}
	}

	if chapterID == 0 {
		return c.JSON(http.StatusNotFound, httpmsgerrorhandler.ErrorResponse{
			Message: "chapter not found",
		})
	}

	// Get pages
	pages, err := h.chapterService.GetPages(c.Request().Context(), chapterID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, pages)
}



func (h Handler) deletePage(c echo.Context) error {
	seriesSlug := c.Param("series_slug")
	chapterNumStr := c.Param("chapter_number")
	pageNumStr := c.Param("page_number")

	// Parse chapter number
	chapterNum, err := strconv.ParseFloat(chapterNumStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter number",
		})
	}

	// Parse page number
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid page number",
		})
	}

	// Get series by full_slug
	series, err := h.seriesService.GetByFullSlug(c.Request().Context(), seriesSlug)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Get all chapters for series to find the chapter ID
	chapters, err := h.chapterService.GetBySeriesID(c.Request().Context(), series.ID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Find chapter by number
	var chapterID uint
	for _, ch := range chapters {
		if ch.ChapterNumber == chapterNum {
			chapterID = ch.ID
			break
		}
	}

	if chapterID == 0 {
		return c.JSON(http.StatusNotFound, httpmsgerrorhandler.ErrorResponse{
			Message: "chapter not found",
		})
	}

	// Get the specific page to get its ID
	page, err := h.chapterService.GetPageByNumber(c.Request().Context(), chapterID, pageNum)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Delete the page
	err = h.chapterService.DeletePage(c.Request().Context(), page.ID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}