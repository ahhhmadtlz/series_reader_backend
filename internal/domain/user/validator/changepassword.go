package validator

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateChangePasswordRequest(ctx context.Context, req param.ChangePasswordRequest) error {
	const op = richerror.Op("uservalidator.ValidateChangePasswordRequest")

	fieldErrors := make(map[string]string)

	err := validation.ValidateStruct(&req,
		validation.Field(&req.OldPassword,
			validation.Required.Error("old password is required"),
		),
		validation.Field(&req.NewPassword,
			validation.Required.Error("new password is required"),
			validation.Length(8, 0).Error("new password must be at least 8 characters long"),
		),
		validation.Field(&req.ConfirmNewPassword,
			validation.Required.Error("confirm password is required"),
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

	// Check if new password and confirm password match
	if req.NewPassword != req.ConfirmNewPassword {
		fieldErrors["confirm_new_password"] = "passwords do not match"
	}

	// Check if new password is different from old password
	if req.OldPassword != "" && req.NewPassword != "" && req.OldPassword == req.NewPassword {
		fieldErrors["new_password"] = "new password must be different from old password"
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