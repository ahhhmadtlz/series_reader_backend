package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) logout(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	var req param.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}
	req.UserID = userID

	if err := h.service.Logout(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "logged out successfully",
	})
}