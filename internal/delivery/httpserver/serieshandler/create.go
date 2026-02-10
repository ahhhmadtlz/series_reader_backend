package serieshandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) create(c echo.Context) error {
	var req param.CreateSeriesRequest

	if err :=c.Bind(&req);err !=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}

	if err := h.validator.ValidateCreateSeriesRequest(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}



	response ,err :=h.service.Create(c.Request().Context(),req)
	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	return c.JSON(http.StatusCreated,response)

}