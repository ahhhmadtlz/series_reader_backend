package user

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) IsPhoneNumberUnique(ctx context.Context, phoneNumber string) (bool, error) {
	const op = richerror.Op("repository.postgres.user.IsPhoneNumberUnique")

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(&exists)
	if err != nil {
		return false, richerror.New(op).WithErr(err).WithMessage("failed to check phone number")
	}

	return !exists, nil
}

func (r *PostgresRepository) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
	const op = richerror.Op("repository.postgres.user.IsUsernameUnique")

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, richerror.New(op).WithErr(err).WithMessage("failed to check username")
	}

	return !exists, nil
}