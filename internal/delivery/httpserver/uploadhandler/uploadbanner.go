package uploadhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) uploadBanner(c echo.Context) error {

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	seriesIDStr := c.Param("seriesID")
	seriesID, err := strconv.ParseUint(seriesIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid series ID",
		})
	}

	file, header, err := c.Request().FormFile("banner")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "failed to read banner file",
		})
	}
	defer file.Close()

	req := param.UploadBannerRequest{
		SeriesID: uint(seriesID),
		UserID:   userID,
		File:     file,
		Header:   header,
	}

	if err := h.validator.ValidateUploadBanner(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	response, err := h.service.UploadBanner(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}