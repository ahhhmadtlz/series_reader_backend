package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) getProfile(c echo.Context) error {

	userID, ok :=c.Get("user_id").(uint)

	if !ok {
		return c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
			Message:"unauthorized",
		})
	}

	response, err:= h.service.GetProfile(c.Request().Context(),userID)

	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	return c.JSON(http.StatusOK, response)

}
