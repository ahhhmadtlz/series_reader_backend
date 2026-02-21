package uploadhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) uploadAvatar(c echo.Context) error {

	userID, ok :=c.Get("user_id").(uint)

	if !ok{
		return c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	file ,header, err :=c.Request().FormFile("avatar")

	if err !=nil{
		return c.JSON(http.StatusBadRequest,httpmsgerrorhandler.ErrorResponse{
			Message: "failed to read avatar file",
		})
	}

	defer file.Close()

	req:=param.UploadAvatarRequest{
		UserID: userID,
		File: file,
		Header: header,
	}

	if err:=h.validator.ValidateUploadAvatar(c.Request().Context(),req); err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	response,err:=h.service.UploadAvatar(c.Request().Context(),req)

	if err !=nil{
		return httpmsgerrorhandler.Error(c,err)
	}

	return c.JSON(http.StatusOK,response)

}