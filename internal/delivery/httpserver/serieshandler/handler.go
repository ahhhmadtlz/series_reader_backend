package serieshandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"
)

type Handler struct {
	service service.Service
	validator validator.Validator
}

func New(service service.Service,validator validator.Validator)Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}