package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) changePassword(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	var req param.ChangePasswordRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}

	// Validate request
	if err := h.validator.ValidateChangePasswordRequest(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	response, err := h.service.ChangePassword(c.Request().Context(), userID, req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}