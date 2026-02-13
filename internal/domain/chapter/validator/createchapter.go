package validator

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateChapterRequest(
	ctx context.Context,
	req param.CreateChapterRequest,
) error {

	const op richerror.Op = "Validator.ValidateCreateChapterRequest"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.SeriesID,
			validation.Required.Error("series_id is required"),
			validation.Min(uint(1)).Error("series_id must be greater than 0"),
		),
		validation.Field(
			&req.ChapterNumber,
			validation.Required.Error("chapter_number is required"),
			validation.Min(0.0).Error("chapter_number must be greater than or equal to 0"),
		),
	)

	if err != nil {
		// Structured validation errors
		if validationErrors, ok := err.(validation.Errors); ok {
			fieldErrors := make(map[string]any, len(validationErrors))

			for field, err := range validationErrors {
				fieldErrors[field] = err.Error()
			}

			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessage("validation failed").
				WithMetaMap(fieldErrors)
		}

		// Unexpected validation failure
		return richerror.Wrap(err, op).
			WithKind(richerror.KindUnexpected).
			WithMessage("unexpected validation error")
	}

	return nil
}
