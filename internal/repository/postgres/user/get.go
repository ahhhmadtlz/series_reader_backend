package user

import (
	"context"
	"database/sql"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetUserByID(ctx context.Context, userID uint) (entity.User, error) {
	const op = richerror.Op("repository.postgres.user.GetUserByID")

	query := `
		SELECT id, username, phone_number, password, avatar_url, bio,
		       role, subscription_tier, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user entity.User
	var roleStr string
	var tierStr string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PhoneNumber,
		&user.Password,
		&user.AvatarURL,
		&user.Bio,
		&roleStr,
		&tierStr,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return entity.User{}, richerror.New(op).
			WithErr(err).
			WithMessage("user not found").
			WithKind(richerror.KindNotFound)
	}

	if err != nil {
		return entity.User{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get user")
	}

	role, _ := sharedentity.MapToRoleEntity(roleStr)
	user.Role = role
	user.SubscriptionTier = sharedentity.MapToSubscriptionTier(tierStr)

	return user, nil
}

func (r *PostgresRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error) {
	const op = richerror.Op("repository.postgres.user.GetUserByPhoneNumber")

	query := `
		SELECT id, username, phone_number, password, avatar_url, bio,
		       role, subscription_tier, is_active, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`

	var user entity.User
	var roleStr string
	var tierStr string

	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.Username,
		&user.PhoneNumber,
		&user.Password,
		&user.AvatarURL,
		&user.Bio,
		&roleStr,
		&tierStr,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return entity.User{}, richerror.New(op).
			WithErr(err).
			WithMessage("user not found").
			WithKind(richerror.KindNotFound)
	}

	if err != nil {
		return entity.User{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get user")
	}

	role, _ := sharedentity.MapToRoleEntity(roleStr)
	user.Role = role
	user.SubscriptionTier = sharedentity.MapToSubscriptionTier(tierStr)

	return user, nil
}