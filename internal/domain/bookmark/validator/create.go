package validator

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateBookmarkRequest(ctx context.Context, req param.CreateBookmarkRequest) error {
	const op = richerror.Op("bookmarkvalidator.ValidateCreateBookmarkRequest")

	fieldErrors := make(map[string]string)

	// Basic validation
	err := validation.ValidateStruct(&req,
		validation.Field(&req.SeriesID,
			validation.Required.Error("series_id is required"),
			validation.Min(uint(1)).Error("series_id must be greater than 0"),
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

	// Return errors if any
	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("invalid input").
			WithKind(richerror.KindInvalid).
			WithMeta("fields", fieldErrors)
	}

	return nil
}