package uploadhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) uploadChapterThumbnail(c echo.Context) error {

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	chapterIDStr := c.Param("chapterID")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	file, header, err := c.Request().FormFile("thumbnail")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "failed to read thumbnail file",
		})
	}
	defer file.Close()

	req := param.UploadChapterThumbnailRequest{
		ChapterID: uint(chapterID),
		UserID:    userID,
		File:      file,
		Header:    header,
	}

	if err := h.validator.ValidateUploadChapterThumbnail(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	response, err := h.service.UploadChapterThumbnail(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}