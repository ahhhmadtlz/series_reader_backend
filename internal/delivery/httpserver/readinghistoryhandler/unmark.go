package readinghistoryhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) unmarkAsRead(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)

	if !ok {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	chapterIDStr := c.Param("chapter_id")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter_id",
		})
	}

	err = h.service.UnmarkAsRead(c.Request().Context(), uint(chapterID), userID)

	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "chapter unmarked as read",
	})
}