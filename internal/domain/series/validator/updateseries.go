package validator

import (
	"context"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateSeriesRequest(
	ctx context.Context,
	req param.UpdateSeriesRequest,
) error {
	const op = richerror.Op("seriesvalidator.ValidateUpdateSeriesRequest")

	fieldErrors := make(map[string]string)

	trim := func(s *string) {
		if s != nil {
			t := strings.TrimSpace(*s)
			*s = t
		}
	}

	trim(req.Title)
	trim(req.Status)
	trim(req.Type)
	trim(req.Description)
	trim(req.Author)
	trim(req.Artist)

	if req.Title != nil {
		if err := validation.Validate(
			req.Title,
			validation.Required.Error("title cannot be empty"),
			validation.Length(1, 255).Error("title must be between 1 and 255 characters"),
		); err != nil {
			fieldErrors["title"] = err.Error()
		}
	}

	if req.Status != nil {
		if err := validation.Validate(
			req.Status,
			validation.In(
				entity.StatusOngoing,
				entity.StatusCompleted,
				entity.StatusHiatus,
				entity.StatusCancelled,
			).Error("status must be one of: ongoing, completed, hiatus, cancelled"),
		); err != nil {
			fieldErrors["status"] = err.Error()
		}
	}

	if req.Type != nil {
		if err := validation.Validate(
			req.Type,
			validation.In(
				entity.TypeManga,
				entity.TypeManhwa,
				entity.TypeManhua,
				entity.TypeComic,
				entity.TypeWebtoon,
			).Error("type must be one of: manga, manhwa, manhua, comic, webtoon"),
		); err != nil {
			fieldErrors["type"] = err.Error()
		}
	}

	if req.Author != nil {
		if err := validation.Validate(
			req.Author,
			validation.Length(0, 255).Error("author must be at most 255 characters"),
		); err != nil {
			fieldErrors["author"] = err.Error()
		}
	}

	if req.Artist != nil {
		if err := validation.Validate(
			req.Artist,
			validation.Length(0, 255).Error("artist must be at most 255 characters"),
		); err != nil {
			fieldErrors["artist"] = err.Error()
		}
	}

	if req.Description != nil {
		if err := validation.Validate(
			req.Description,
			validation.Length(0, 5000).Error("description must be at most 5000 characters"),
		); err != nil {
			fieldErrors["description"] = err.Error()
		}
	}

	if req.CoverImageURL != nil {
		if err := validation.Validate(
			req.CoverImageURL,
			validation.Length(0, 500).Error("cover image URL must be at most 500 characters"),
		); err != nil {
			fieldErrors["cover_image_url"] = err.Error()
		}
	}

	if req.PublicationYear != nil {
		if err := validation.Validate(
			req.PublicationYear,
			validation.Min(1900).Error("publication year must be 1900 or later"),
			validation.Max(2150).Error("publication year must be 2150 or earlier"),
		); err != nil {
			fieldErrors["publication_year"] = err.Error()
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
