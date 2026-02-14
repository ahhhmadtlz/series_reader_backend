package validator

import (
	"context"
	"errors"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateLoginRequest(ctx context.Context, req param.LoginRequest) error {
	const op = richerror.Op("uservalidator.ValidateLoginRequest")

	// Trim spaces
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)

	fieldErrors := make(map[string]string)

	// Basic validation
	err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("phone number is required"),
			validation.Match(phoneNumberRegex).Error("phone number format is invalid (should be 09xxxxxxxxx)"),
		),
		validation.Field(&req.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 0).Error("password must be at least 8 characters long"),
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

	// Check if phone number exists (only if format is valid)
	if req.PhoneNumber != "" && phoneNumberRegex.MatchString(req.PhoneNumber) {
		_, err := v.repo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
		if err != nil {
			// Check if it's a "not found" error
			var re *richerror.RichError
			if errors.As(err, &re) && re.GetKind() == richerror.KindNotFound {
				fieldErrors["phone_number"] = "phone number not found"
			} else {
				// Unexpected database error
				return richerror.New(op).
					WithMessage("failed to check phone number").
					WithKind(richerror.KindUnexpected).
					WithErr(err)
			}
		}
	}

	// Return error with embedded field errors
	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("invalid input").
			WithKind(richerror.KindInvalid).
			WithMeta("fields", fieldErrors)
	}

	return nil
}