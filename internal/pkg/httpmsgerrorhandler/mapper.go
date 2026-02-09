package httpmsgerrorhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string         `json:"message"`
	Errors  map[string]any `json:"errors,omitempty"`
}

func Error(c echo.Context, err error) error {
	re,ok :=err.(*richerror.RichError)

	if !ok{
		return c.JSON(http.StatusInternalServerError,ErrorResponse{
			Message: "internal server error",
		})
	}
	statusCode :=mapKindToHTTPStatus(re.GetKind())
	message:=re.GetMessage()
	if message== ""{
		message= re.Kind.String()
	}

	response :=ErrorResponse{
		Message: message,
	}

	if re.Meta !=nil{
		response.Errors=re.Meta
	}
	return c.JSON(statusCode,response)
}


func mapKindToHTTPStatus(kind richerror.Kind)int {
	switch kind {
		case richerror.KindInvalid:
			return http.StatusBadRequest
		case richerror.KindForbidden:
			return http.StatusForbidden
		case richerror.KindNotFound:
			return http.StatusNotFound
		case richerror.KindUnexpected:
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
	}
}