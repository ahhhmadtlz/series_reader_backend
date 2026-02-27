package user

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) SaveRefreshToken(ctx context.Context, userID uint, tokenHash string, expiresAt time.Time) error {
	const op = richerror.Op("repository.postgres.user.SaveRefreshToken")

	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to save refresh token").WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (r *PostgresRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	const op = richerror.Op("repository.postgres.user.RevokeRefreshToken")

	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to revoke refresh token").WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (r *PostgresRepository) IsRefreshTokenValid(ctx context.Context, tokenHash string) (bool, error) {
	const op = richerror.Op("repository.postgres.user.IsRefreshTokenValid")

	query := `
		SELECT EXISTS(
			SELECT 1 FROM refresh_tokens
			WHERE token_hash = $1
			  AND revoked_at IS NULL
			  AND expires_at > NOW()
		)
	`

	var valid bool
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&valid)
	if err != nil {
		return false, richerror.New(op).WithErr(err).WithMessage("failed to check refresh token").WithKind(richerror.KindUnexpected)
	}

	return valid, nil
}