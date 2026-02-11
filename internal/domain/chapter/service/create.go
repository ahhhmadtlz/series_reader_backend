package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) Create(ctx context.Context, req param.CreateChapterRequest) (param.ChapterResponse, error) {
	const op = richerror.Op("service.chapter.Create")

	chapter := &entity.Chapter{
		SeriesID:      req.SeriesID,
		ChapterNumber: req.ChapterNumber,
		Title:         req.Title,
	}

	created, err := s.repo.Create(ctx, chapter)
	if err != nil {
		return param.ChapterResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to create chapter").WithKind(richerror.KindUnexpected)
	}

	return toChapterResponse(created), nil
}