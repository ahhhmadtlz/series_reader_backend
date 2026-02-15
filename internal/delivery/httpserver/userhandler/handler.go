package userhandler

import (
	userservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/validator"
)

type Handler struct {
	service   userservice.Service
	validator uservalidator.Validator
}

func New(service userservice.Service, validator uservalidator.Validator) Handler {
	return Handler{
		service:   service,
		validator: validator,
	}
}