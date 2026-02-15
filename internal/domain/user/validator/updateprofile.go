package validator

import (
	"context"
	"strings"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateProfileRequest(ctx context.Context, userID uint, req param.UpdateProfileRequest) error {
	const op = richerror.Op("uservalidator.ValidateUpdateProfileRequest")

	// Trim spaces
	req.Username = strings.TrimSpace(req.Username)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)
	req.Bio = strings.TrimSpace(req.Bio)

	fieldErrors := make(map[string]string)

	// Basic validation
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Required.Error("username is required"),
			validation.Length(3, 30).Error("username must be between 3 and 30 characters"),
			validation.Match(usernameRegex).Error("username can only contain letters, numbers, underscore and hyphen"),
		),
		validation.Field(&req.AvatarURL,
			validation.When(req.AvatarURL != "",
				validation.Match(urlRegex).Error("avatar URL must be a valid URL starting with http:// or https://"),
			),
		),
		validation.Field(&req.Bio,
			validation.Length(0, 500).Error("bio must be at most 500 characters"),
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

	// Get current user to check username change
	currentUser, err := v.repo.GetUserByID(ctx, userID)
	if err != nil {
		return richerror.New(op).
			WithMessage("failed to get user").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	// Check if username is being changed
	if req.Username != currentUser.Username {
		// Check if username was changed in the last 12 months
		if currentUser.UsernameLastChangedAt != nil {
			twelveMonthsAgo := time.Now().AddDate(0, -12, 0)
			if currentUser.UsernameLastChangedAt.After(twelveMonthsAgo) {
				nextChangeDate := currentUser.UsernameLastChangedAt.AddDate(0, 12, 0)
				fieldErrors["username"] = "username can only be changed once every 12 months. Next change available on: " + nextChangeDate.Format("2006-01-02")
			}
		}

		// Check username uniqueness (only if format is valid and 12-month rule passed)
		if req.Username != "" && usernameRegex.MatchString(req.Username) && fieldErrors["username"] == "" {
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