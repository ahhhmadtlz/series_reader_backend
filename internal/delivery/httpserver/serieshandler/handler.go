package serieshandler

import (
	seriesService "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"
)

type Handler struct {
	service seriesService.Service
	validator validator.Validator
}

func New(service seriesService.Service,validator validator.Validator)Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}