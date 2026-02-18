package adminhandler

import (
	userService "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/service"
)

type Handler struct {
	userService userService.Service
}

func New(userService userService.Service) Handler {
	return Handler{
		userService: userService,
	}
}