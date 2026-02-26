package chapterhandler

import (
	"context"

	chapterparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	chapterValidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/validator"
	seriesparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
)

type ChapterService interface {
    Create(ctx context.Context, req chapterparam.CreateChapterRequest) (chapterparam.ChapterResponse, error)
    GetBySeriesID(ctx context.Context, seriesID uint) ([]chapterparam.ChapterResponse, error)
    GetChapterWithPages(ctx context.Context, chapterID uint) (chapterparam.ChapterWithPagesResponse, error)
    Delete(ctx context.Context, id uint) error
    UploadPage(ctx context.Context, req chapterparam.UploadPageParam) (chapterparam.ChapterPageResponse, error)
    BulkUploadPages(ctx context.Context, req chapterparam.BulkUploadParam) ([]chapterparam.ChapterPageResponse, error)
		ReorderPages(ctx context.Context, req chapterparam.ReorderPagesParam) error
    GetPages(ctx context.Context, chapterID uint) ([]chapterparam.ChapterPageResponse, error)
    DeletePage(ctx context.Context, chapterID uint, pageNumber int) error
}

type SeriesService interface {
    GetByFullSlug(ctx context.Context, slug string) (seriesparam.SeriesResponse, error)
}

type Handler struct {
    chapterService   ChapterService
    seriesService    SeriesService
    chapterValidator chapterValidator.Validator
}

func New(chapterService ChapterService, seriesService SeriesService, chapterValidator chapterValidator.Validator) Handler {
    return Handler{
        chapterService:   chapterService,
        seriesService:    seriesService,
        chapterValidator: chapterValidator,
    }
}