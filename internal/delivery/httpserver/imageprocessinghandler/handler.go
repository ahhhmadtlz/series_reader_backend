package imageprocessinghandler

import (
	ipservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/service"
)

type Handler struct {
	ipService ipservice.Service
}

func New(ipService ipservice.Service) Handler {
	return Handler{
		ipService: ipService,
	}
}