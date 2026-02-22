package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) ReorderPages(ctx context.Context, req param.ReorderPagesParam) error {
	const op = richerror.Op("service.chapter.ReorderPages")

	// 1. Verify chapter exists
	_, err := s.repo.GetByID(ctx, req.ChapterID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("chapter not found").
			WithKind(richerror.KindNotFound)
	}

	if len(req.Pages) == 0 {
		return richerror.New(op).
			WithMessage("no pages provided").
			WithKind(richerror.KindInvalid)
	}

	// 2. Build updates
	updates := make([]entity.PageNumberUpdate, len(req.Pages))
	for i, p := range req.Pages {
		if p.PageNumber < 0 {
			return richerror.New(op).
				WithMessage("page number must be non-negative").
				WithKind(richerror.KindInvalid)
		}
		updates[i] = entity.PageNumberUpdate{
			PageID:     p.PageID,
			PageNumber: p.PageNumber,
		}
	}

	// 3. Update in DB (uses transaction in repo)
	if err := s.repo.UpdatePageNumbers(ctx, updates); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to reorder pages").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("pages reordered", "chapter_id", req.ChapterID, "count", len(updates))

	return nil
}