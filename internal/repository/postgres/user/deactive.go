package user

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) DeactivateUser(ctx context.Context,userID uint)error {
	const op= richerror.Op("repository.postgres.user.DeactivateUser")

	query:=`UPDATE users SET is_active = false WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)

	if err !=nil {
		return  richerror.New(op).WithErr(err).WithMessage("failed to deactivate user")
	}

	rowsAffected,err :=result.RowsAffected()
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).WithMessage("user not found").WithKind(richerror.KindNotFound)
	}

	return nil
}