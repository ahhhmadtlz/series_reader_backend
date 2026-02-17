package readinghistoryhandler

import (
	readinghistoryService "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/validator"
)

type Handler struct {
	service   readinghistoryService.Service
	validator validator.Validator
}

func New(service readinghistoryService.Service, validator validator.Validator) Handler {
	return Handler{
		service:   service,
		validator: validator,
	}
}