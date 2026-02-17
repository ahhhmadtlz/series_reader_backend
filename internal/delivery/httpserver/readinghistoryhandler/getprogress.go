package readinghistoryhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) getSeriesProgress(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)

	if !ok {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	seriesSlug := c.Param("slug")

	if seriesSlug == "" {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "series slug is required",
		})
	}

	response, err := h.service.GetSeriesProgress(c.Request().Context(), userID, seriesSlug)

	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}