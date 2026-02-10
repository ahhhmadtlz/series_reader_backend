package serieshandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

// update handles PUT /series/:id
func (h Handler) update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid series ID",
		})
	}

	var req param.UpdateSeriesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}

	// Validate request
	if err := h.validator.ValidateUpdateSeriesRequest(
		c.Request().Context(),
		req,
	); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Update series
	response, err := h.service.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}