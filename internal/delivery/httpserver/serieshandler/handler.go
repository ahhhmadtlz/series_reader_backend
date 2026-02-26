package serieshandler

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/validator"
)

//why am i using series service here it will apply to all handler
// Think of it like this:
// The interface says "I need a car that has 4 wheels and starts"
// The concrete service is "here is a Toyota"
// The handler doesn't need to know it's a Toyota
type SeriesService interface {
    Create(ctx context.Context, req param.CreateSeriesRequest) (param.SeriesResponse, error)
    GetByID(ctx context.Context, id uint) (param.SeriesResponse, error)
    GetByFullSlug(ctx context.Context, slug string) (param.SeriesResponse, error)
    GetList(ctx context.Context, req param.GetListRequest) (param.GetListResponse, error)
    Update(ctx context.Context, id uint, req param.UpdateSeriesRequest) (param.SeriesResponse, error)
    Delete(ctx context.Context, id uint) error
    IncrementViewCount(ctx context.Context, id uint) error
}
type Handler struct {
    service   SeriesService
    validator validator.Validator
}

func New(service SeriesService, validator validator.Validator) Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}