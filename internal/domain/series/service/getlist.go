package service

import (
	"context"
	"math"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetList(ctx context.Context, req param.GetListRequest) (param.GetListResponse, error) {
	const op = richerror.Op("service.series.GetList")

	// Set default pagination
	if req.Pagination.Page < 1 {
		req.Pagination.Page = 1
	}
	if req.Pagination.PageSize < 1 || req.Pagination.PageSize > 100 {
		req.Pagination.PageSize = 20
	}

	// Set default sorting
	if req.Sort.SortBy == "" {
		req.Sort.SortBy = "created_at"
	}
	if req.Sort.SortOrder == "" {
		req.Sort.SortOrder = "desc"
	}

	seriesList, totalCount, err := s.repo.GetList(ctx, req)
	if err != nil {
		return param.GetListResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to get series list").WithKind(richerror.KindUnexpected)
	}

	// Convert entities to responses
	items := make([]param.SeriesResponse, len(seriesList))
	for i, series := range seriesList {
		items[i] = toSeriesResponse(series)
	}

	totalPages := int(math.Ceil(
		float64(totalCount) / float64(req.Pagination.PageSize),
	))

	return param.GetListResponse{
		Items: items,
		Pagination: param.PaginationMeta{
			Page:       req.Pagination.Page,
			PageSize:   req.Pagination.PageSize,
			TotalItems: totalCount,
			TotalPages: totalPages,
		},
	}, nil
}