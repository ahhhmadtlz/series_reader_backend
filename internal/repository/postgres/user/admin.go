package user

import (
	"context"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) UpdateUserRole(ctx context.Context, userID uint, role sharedentity.Role) error {
	const op = richerror.Op("repository.postgres.user.UpdateUserRole")

	query := `
		UPDATE users
		SET role = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, role.String(), userID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to update user role")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("user not found").
			WithKind(richerror.KindNotFound)
	}

	return nil
}

func (r *PostgresRepository) GrantPermission(ctx context.Context, userID uint, permission string, grantedBy uint) error {
	const op = richerror.Op("repository.postgres.user.GrantPermission")

	query := `
	 INSERT INTO user_permissions(user_id, permission,granted_by)
	 VALUES ($1, $2, $3)
	 ON CONFLICT (user_id,permission)
	 DO UPDATE SET granted_by = $3, granted_at = NOW()
	`
	_,err:=r.db.ExecContext(ctx,query,userID,permission,grantedBy)
	if err !=nil{
		return richerror.New(op).
		  WithErr(err).
			WithMessage("failed to grant permission")
	}
	return  nil
}

func (r *PostgresRepository) RevokePermission(ctx context.Context, userID uint, permission string) error {
	const op = richerror.Op("repository.postgres.user.RevokePermission")

	query := `
		DELETE FROM user_permissions
		WHERE user_id = $1 AND permission = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, permission)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to revoke permission")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("permission not found").
			WithKind(richerror.KindNotFound)
	}

	return nil
}

