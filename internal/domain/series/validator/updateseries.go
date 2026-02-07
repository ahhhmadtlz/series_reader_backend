package validator

import (
	"context"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateSeriesRequest(ctx context.Context, req param.UpdateSeriesRequest) (map[string]string, error) {
	const op = richerror.Op("seriesvalidator.ValidateUpdateSeriesRequest")
	
	fieldErrors:=make(map[string]string)

	if req.Title !=nil{
		trimmed:=strings.TrimSpace(*req.Title)
		req.Title=&trimmed
	}

	if req.Status !=nil{
		trimmed:=strings.TrimSpace(*req.Status)
		req.Status=&trimmed
	}
	
	if req.Type != nil {
		trimmed := strings.TrimSpace(*req.Type)
		req.Type = &trimmed
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		req.Description = &trimmed
	}
	if req.Author != nil {
		trimmed := strings.TrimSpace(*req.Author)
		req.Author = &trimmed
	}
	if req.Artist != nil {
		trimmed := strings.TrimSpace(*req.Artist)
		req.Artist = &trimmed
	}
	if req.Title !=nil{
		err:=validation.Validate(req.Title,
			validation.Required.Error("title cannot be empty"),
			validation.Length(1, 255).Error("title must be between 1 and 255 characters"),
		)
		if err != nil {
			fieldErrors["title"] = err.Error()
		}
	}
	if req.Status != nil {
		err := validation.Validate(req.Status,
			validation.In(
				entity.StatusOngoing,
				entity.StatusCompleted,
				entity.StatusHiatus,
				entity.StatusCancelled,
			).Error("status must be one of: ongoing, completed, hiatus, cancelled"),
		)
		if err != nil {
			fieldErrors["status"] = err.Error()
		}
	}

	// Validate Type if provided
	if req.Type != nil {
		err := validation.Validate(req.Type,
			validation.In(
				entity.TypeManga,
				entity.TypeManhwa,
				entity.TypeManhua,
				entity.TypeComic,
				entity.TypeWebtoon,
			).Error("type must be one of: manga, manhwa, manhua, comic, webtoon"),
		)
		if err != nil {
			fieldErrors["type"] = err.Error()
		}
	}

	// Validate Author if provided
	if req.Author != nil {
		err := validation.Validate(req.Author,
			validation.Length(0, 255).Error("author must be at most 255 characters"),
		)
		if err != nil {
			fieldErrors["author"] = err.Error()
		}
	}

	// Validate Artist if provided
	if req.Artist != nil {
		err := validation.Validate(req.Artist,
			validation.Length(0, 255).Error("artist must be at most 255 characters"),
		)
		if err != nil {
			fieldErrors["artist"] = err.Error()
		}
	}

	// Validate Description if provided
	if req.Description != nil {
		err := validation.Validate(req.Description,
			validation.Length(0, 5000).Error("description must be at most 5000 characters"),
		)
		if err != nil {
			fieldErrors["description"] = err.Error()
		}
	}

	// Validate CoverImageURL if provided
	if req.CoverImageURL != nil {
		err := validation.Validate(req.CoverImageURL,
			validation.Length(0, 500).Error("cover image URL must be at most 500 characters"),
		)
		if err != nil {
			fieldErrors["cover_image_url"] = err.Error()
		}
	}

	// Validate PublicationYear if provided
	if req.PublicationYear != nil {
		err := validation.Validate(req.PublicationYear,
			validation.Min(1900).Error("publication year must be 1900 or later"),
			validation.Max(2150).Error("publication year must be 2150 or earlier"),
		)
		if err != nil {
			fieldErrors["publication_year"] = err.Error()
		}
	}

	if len(fieldErrors) > 0 {
		return fieldErrors, richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields", fieldErrors)
	}

	return nil, nil

}