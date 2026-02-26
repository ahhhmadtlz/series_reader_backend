package readinghistoryhandler

import (
	"context"

	rhparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/validator"
)

type ReadingHistoryService interface {
    MarkAsRead(ctx context.Context, req rhparam.MarkAsReadRequest, userID uint) (rhparam.ReadingHistoryResponse, error)
    UnmarkAsRead(ctx context.Context, chapterID uint, userID uint) error
    GetUserHistory(ctx context.Context, userID uint, limit int, offset int) ([]rhparam.ReadingHistoryResponse, error)
    GetSeriesProgress(ctx context.Context, userID uint, seriesSlug string) (rhparam.SeriesProgressResponse, error)
}

type Handler struct {
    service   ReadingHistoryService
    validator validator.Validator
}

func New(service ReadingHistoryService, validator validator.Validator) Handler {
    return Handler{
        service:   service,
        validator: validator,
    }
}