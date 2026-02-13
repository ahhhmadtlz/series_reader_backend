package validator

import (
	"context"
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddChapterPagesRequest(
	ctx context.Context,
	req param.AddChapterPagesRequest,
) error {

	const op richerror.Op = "Validator.ValidateAddChapterPagesRequest"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.ChapterID,
			validation.Required.Error("chapter_id is required"),
			validation.Min(uint(1)).Error("chapter_id must be greater than 0"),
		),
		validation.Field(
			&req.Pages,
			validation.Required.Error("pages is required"),
			validation.Length(1, 0).Error("at least one page is required"),
		),
	)

	if err != nil {
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

		return richerror.Wrap(err, op).
			WithKind(richerror.KindUnexpected).
			WithMessage("unexpected validation error")
	}

	// Validate each page
	fieldErrors := make(map[string]any)

	for i, page := range req.Pages {
		if err := validatePageItem(page, i); err != nil {

			if validationErrors, ok := err.(validation.Errors); ok {
				for field, err := range validationErrors {
					key := fmt.Sprintf("pages[%d].%s", i, field)
					fieldErrors[key] = err.Error()
				}
				continue
			}

			return richerror.Wrap(err, op).
				WithKind(richerror.KindUnexpected).
				WithMessage("unexpected page validation error")
		}
	}

	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessage("validation failed").
			WithMetaMap(fieldErrors)
	}

	return nil
}



func validatePageItem(page param.CreateChapterPageItem, _ int) error {
	return validation.ValidateStruct(&page,
		validation.Field(&page.PageNumber,
			validation.Required.Error("page_number is required"),
			validation.Min(1).Error("page_number must be greater than 0"),
		),
		validation.Field(&page.ImageURL,
			validation.Required.Error("image_url is required"),
			validation.NotNil.Error("image_url cannot be empty"),
		),
	)
}