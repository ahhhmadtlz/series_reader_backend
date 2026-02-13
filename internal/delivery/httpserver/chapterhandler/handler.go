package chapterhandler

import (
	chapterService "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/service"
	chapterValidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/validator"
	seriesService "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/service"
)

type Handler struct {
	chapterService chapterService.Service
	seriesService seriesService.Service
	chapterValidator  chapterValidator.Validator
}


func New(chapterService chapterService.Service,seriesService seriesService.Service,chapterValidator chapterValidator.Validator)Handler{
	return Handler{
		chapterService: chapterService,
		seriesService: seriesService,
		chapterValidator: chapterValidator,
	}
}

