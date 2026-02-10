package serieshandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

// delete handles DELETE /series/:id
func (h Handler) delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid series ID",
		})
	}

	err = h.service.Delete(c.Request().Context(), uint(id))
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}