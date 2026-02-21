package uploadhandler

import (
	uploadservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/service"
	uploadvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/validator"
)

type Handler struct {
	service uploadservice.Service
	validator uploadvalidator.Validator
}


func New(service uploadservice.Service, validator uploadvalidator.Validator) Handler {
	return  Handler{
		service: service,
		validator: validator ,
	}
}