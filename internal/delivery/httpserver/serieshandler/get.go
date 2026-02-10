package serieshandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) get(c echo.Context) error {
    identifier := c.Param("identifier")

    // 1. Try to parse as an Integer (ID)
    id, err := strconv.ParseUint(identifier, 10, 32)
    
    if err == nil {
        // SUCCESS: It was a number. Call GetByID logic.
        response, err := h.service.GetByID(c.Request().Context(), uint(id))
        if err != nil {
            return httpmsgerrorhandler.Error(c, err)
        }
        _ = h.service.IncrementViewCount(c.Request().Context(), uint(id))
        return c.JSON(http.StatusOK, response)
    }

    // 2. Failure: It was not a number. Treat it as a FullSlug.
    // Ensure you call GetByFullSlug here (matching the service interface)
    response, err := h.service.GetByFullSlug(c.Request().Context(), identifier)
    if err != nil {
        // Use your existing error handling logic
        if re, ok := err.(*richerror.RichError); ok && re.GetKind() == richerror.KindNotFound {
            return c.JSON(http.StatusNotFound, httpmsgerrorhandler.ErrorResponse{
                Message: "series not found",
            })
        }
        return httpmsgerrorhandler.Error(c, err)
    }

   _ = h.service.IncrementViewCount(c.Request().Context(), response.ID)
    return c.JSON(http.StatusOK, response)
}