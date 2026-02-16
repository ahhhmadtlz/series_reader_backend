package bookmarkhandler

import (
	bookmarkService "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/validator"
)

type Handler struct {
	service  bookmarkService.Service
	validator validator.Validator
}


func New(service bookmarkService.Service, validator validator.Validator)Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}
