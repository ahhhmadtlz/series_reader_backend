package validator

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateGetListRequest(
	ctx context.Context,
	req param.GetListRequest,
) error {
	const op = richerror.Op("seriesvalidator.ValidateGetListRequest")

	fieldErrors := make(map[string]string)

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Filter.Status,
			validation.In(
				"",
				entity.StatusOngoing,
				entity.StatusCompleted,
				entity.StatusHiatus,
				entity.StatusCancelled,
			).Error("invalid status filter"),
		),
		validation.Field(&req.Filter.Type,
			validation.In(
				"",
				entity.TypeManga,
				entity.TypeManhwa,
				entity.TypeManhua,
				entity.TypeComic,
				entity.TypeWebtoon,
			).Error("invalid type filter"),
		),
		validation.Field(&req.Sort.SortBy,
			validation.In(
				"",
				"rating",
				"view_count",
				"created_at",
				"title",
			).Error("sort_by must be one of: rating, view_count, created_at, title"),
		),
		validation.Field(&req.Sort.SortOrder,
			validation.In("", "asc", "desc").
				Error("sort_order must be asc or desc"),
		),
		validation.Field(&req.Pagination.Page,
			validation.Min(0).Error("page must be a positive number"),
		),
		validation.Field(&req.Pagination.PageSize,
			validation.Min(0).Error("page_size must be a positive number"),
			validation.Max(100).Error("page_size must be at most 100"),
		),
	)

	if err != nil {
		if errV, ok := err.(validation.Errors); ok {
			for field, value := range errV {
				if value != nil {
					fieldErrors[field] = value.Error()
				}
			}
		}
	}

	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("invalid input").
			WithKind(richerror.KindInvalid).
			WithMeta("fields", fieldErrors)
	}

	return nil
}
