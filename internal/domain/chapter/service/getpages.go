package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetPages(ctx context.Context, chapterID uint) ([]param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.GetPages")

	pages, err := s.repo.GetPagesByChapterID(ctx, chapterID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get pages").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.ChapterPageResponse, len(pages))
	for i, p := range pages {
		responses[i] = param.ChapterPageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
		}
	}

	return responses, nil
}

func (s Service) GetPageByNumber(ctx context.Context, chapterID uint, pageNumber int) (param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.GetPageByNumber")

	page, err := s.repo.GetPageByNumber(ctx, chapterID, pageNumber)
	if err != nil {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("page not found").
			WithKind(richerror.KindNotFound)
	}

	return param.ChapterPageResponse{
		ID:         page.ID,
		PageNumber: page.PageNumber,
		ImageURL:   page.ImageURL,
	}, nil
}