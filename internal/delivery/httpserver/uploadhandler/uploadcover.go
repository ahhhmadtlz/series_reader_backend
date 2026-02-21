package uploadhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) uploadCover(c echo.Context) error {

	userID, ok := c.Get("user_id").(uint)

	if !ok {
		return c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	seriesIDStr := c.Param("id")
	seriesID,err :=strconv.ParseUint(seriesIDStr,10,32)

	if err !=nil{
		return  c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "invalid series ID",
		})
	}

	file , header, err :=c.Request().FormFile("cover")

	if err !=nil {
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "failed to read cover file",
		})
	}

	defer file.Close()

	req:=param.UploadCoverRequest{
		SeriesID: uint(seriesID),
		UserID: userID,
		File: file,
		Header: header,
	}

		// 5. Validate
	if err := h.validator.ValidateUploadCover(c.Request().Context(), req); err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	// 6. Call service
	response, err := h.service.UploadCover(c.Request().Context(), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)

}