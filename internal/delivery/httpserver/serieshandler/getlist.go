package serieshandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

// getList handles GET /series with query parameters for filtering, sorting, and pagination
func (h Handler) getList(c echo.Context) error {
	var req param.GetListRequest

	// Bind filter parameters
	if err := c.Bind(&req.Filter); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid filter parameters",
		})
	}

	// Bind sort parameters
	if err := c.Bind(&req.Sort); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid sort parameters",
		})
	}

	// Bind pagination parameters
	if err := c.Bind(&req.Pagination); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid pagination parameters",
		})
	}

	// Validate request
	 err := h.validator.ValidateGetListRequest(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// Get series list
	response, err := h.service.GetList(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}