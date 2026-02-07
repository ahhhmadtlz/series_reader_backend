package validator

import (
	"context"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateSeriesRequest(ctx context.Context, req param.CreateSeriesRequest) (map[string]string, error) {
	const op = richerror.Op("seriesvalidator.ValidateCreateSeriesRequest")

	// Trim spaces
	req.Title = strings.TrimSpace(req.Title)
	req.Status = strings.TrimSpace(req.Status)
	req.Type = strings.TrimSpace(req.Type)
	req.Description = strings.TrimSpace(req.Description)
	req.Author = strings.TrimSpace(req.Author)
	req.Artist = strings.TrimSpace(req.Artist)

	fieldErrors := make(map[string]string)

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.Required.Error("title is required"),
			validation.Length(1, 255).Error("title must be between 1 and 255 characters"),
		),
		validation.Field(&req.Status,
			validation.Required.Error("status is required"),
			validation.In(
				entity.StatusOngoing,
				entity.StatusCompleted,
				entity.StatusHiatus,
				entity.StatusCancelled,
			).Error("status must be one of: ongoing, completed, hiatus, cancelled"),
		),
		validation.Field(&req.Type,
			validation.Required.Error("type is required"),
			validation.In(
				entity.TypeManga,
				entity.TypeManhwa,
				entity.TypeManhua,
				entity.TypeComic,
				entity.TypeWebtoon,
			).Error("type must be one of: manga, manhwa, manhua, comic, webtoon"),
		),
		validation.Field(&req.Author,
			validation.Length(0, 255).Error("author must be at most 255 characters"),
		),
		validation.Field(&req.Artist,
			validation.Length(1, 255).Error("artist must be at most 255 characters"),
		),
		validation.Field(&req.Description,
			validation.Length(0, 5000).Error("description must be at most 5000 characters"),
		),
		validation.Field(&req.CoverImageURL,
			validation.Length(0, 500).Error("cover image URL must be at most 500 characters"),
		),
		validation.Field(&req.PublicationYear,
			validation.Min(1900).Error("publication year must be 1900 or later"),
			validation.Max(2100).Error("publication year must be 2100 or earlier"),
		),
	)

	if err != nil {
		if errV, ok := err.(validation.Errors); ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}
	}

	if len(fieldErrors) > 0 {
		return fieldErrors, richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields", fieldErrors)
	}

	return nil, nil
}