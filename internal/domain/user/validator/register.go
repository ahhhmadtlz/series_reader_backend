package validator

import (
	"context"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateRegisterRequest(ctx context.Context, req param.RegisterRequest) error {
	const op = richerror.Op("uservalidator.ValidateRegisterRequest")

	// Trim spaces
	req.Username = strings.TrimSpace(req.Username)
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)

	fieldErrors := make(map[string]string)

	// Basic validation
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Required.Error("username is required"),
			validation.Length(3, 30).Error("username must be between 3 and 30 characters"),
			validation.Match(usernameRegex).Error("username can only contain letters, numbers, underscore and hyphen"),
		),
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

	// Check phone number uniqueness (only if format is valid)
	if req.PhoneNumber != "" && phoneNumberRegex.MatchString(req.PhoneNumber) {
		isUnique, err := v.repo.IsPhoneNumberUnique(ctx, req.PhoneNumber)
		if err != nil {
			return richerror.New(op).
				WithMessage("failed to check phone number").
				WithKind(richerror.KindUnexpected).
				WithErr(err)
		}
		if !isUnique {
			fieldErrors["phone_number"] = "phone number is already registered"
		}
	}

	// Check username uniqueness (only if format is valid)
	if req.Username != "" && usernameRegex.MatchString(req.Username) {
		isUnique, err := v.repo.IsUsernameUnique(ctx, req.Username)
		if err != nil {
			return richerror.New(op).
				WithMessage("failed to check username").
				WithKind(richerror.KindUnexpected).
				WithErr(err)
		}
		if !isUnique {
			fieldErrors["username"] = "username is already taken"
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