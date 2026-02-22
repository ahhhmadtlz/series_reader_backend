package chapterhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

//uploadPage handles POST /chapters/:id/pages

func (h Handler) uploadPage(c echo.Context) error {
	chapterID, err:=parseID(c,"id")
	if err !=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	pageNumberStr :=c.FormValue("page_number")
	pageNumber,err :=strconv.Atoi(pageNumberStr)

	if err !=nil {
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid page number",
		})
	}

	file,header,err :=c.Request().FormFile("page")
	if err !=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "failed to read page file",
		})
	}
	defer file.Close()

	req:=param.UploadPageParam{
		ChapterID: chapterID,
		PageNumber: pageNumber,
		File: file,
		Header: header,
	}
	if err := h.chapterValidator.ValidateUploadPage(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	response, err := h.chapterService.UploadPage(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusCreated, response)

}

// bulkUploadPages handles POST /chapters/:id/pages/bulk
func (h Handler) bulkUploadPages(c echo.Context) error {
	chapterID, err := parseID(c, "id")

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "failed to parse multipart form",
		})
	}

	files := form.File["pages"]

	req := param.BulkUploadParam{
		ChapterID: chapterID,
		Files:     files,
	}

	if err := h.chapterValidator.ValidateBulkUpload(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	responses, err := h.chapterService.BulkUploadPages(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "pages uploaded successfully",
		"pages":   responses,
	})
}

// reorderPages handles PATCH /chapters/:id/pages/reorder
func (h Handler) reorderPages(c echo.Context) error {
	chapterID, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	var req param.ReorderPagesParam
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}
	req.ChapterID = chapterID

	if err := h.chapterValidator.ValidateReorderPages(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	if err := h.chapterService.ReorderPages(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "pages reordered successfully",
	})
}

// getPages handles GET /chapters/:id/pages
func (h Handler) getPages(c echo.Context) error {
	chapterID, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	pages, err := h.chapterService.GetPages(c.Request().Context(), chapterID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, pages)
}

// deletePage handles DELETE /chapters/:id/pages/:page_number
func (h Handler) deletePage(c echo.Context) error {
	chapterID, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	pageNumber, err := strconv.Atoi(c.Param("page_number"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid page number",
		})
	}

	if err := h.chapterService.DeletePage(c.Request().Context(), chapterID, pageNumber); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// parseID is a helper to parse uint path params
func parseID(c echo.Context, paramName string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(paramName), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}