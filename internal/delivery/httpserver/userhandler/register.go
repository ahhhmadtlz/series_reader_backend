package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) register(c echo.Context) error {
	var req param.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}

	if err := h.validator.ValidateRegisterRequest(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	response, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusCreated, response)
}