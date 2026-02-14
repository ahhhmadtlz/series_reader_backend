package validator

import (
	"context"
	"regexp"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
)

var (
	// Phone number regex for Iranian format: 09xxxxxxxxx (11 digits)
	phoneNumberRegex = regexp.MustCompile(`^09[0-9]{9}$`)
	
	// Username: alphanumeric, underscore, hyphen (3-30 chars)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)
)

type Repository interface {
	IsPhoneNumberUnique(ctx context.Context, phoneNumber string) (bool, error)
	IsUsernameUnique(ctx context.Context, username string) (bool, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}