package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) updateProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)

	if !ok {
		return c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	var req param.UpdateProfileRequest

	if err:=c.Bind(&req);err!=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}
	
	if err := h.validator.ValidateUpdateProfileRequest(c.Request().Context(), userID, req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	response, err:=h.service.UpdateProfile(c.Request().Context(),userID,req)

	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	return c.JSON(http.StatusOK,response)
}