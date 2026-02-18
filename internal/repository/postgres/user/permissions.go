package user

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	const op = richerror.Op("repository.postgres.user.GetUserPermissions")
	
	query := `
		SELECT permission
		FROM user_permissions
		WHERE user_id = $1
	`
  rows, err :=r.db.QueryContext(ctx, query, userID)

	if err !=nil{
		return nil, richerror.New(op).
		 WithErr(err).
		 WithMessage("failed to get user permissions")
	}
	defer rows.Close()

	var permissions []string

	for rows.Next(){
		var permission string
		if err :=rows.Scan(&permission);err !=nil{
			return nil, richerror.New(op).
			 WithErr(err).
			 WithMessage("failed to scan permission")
		}
		permissions =append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("error iterating permissions")
	}

	return permissions, nil
}

func (r *PostgresRepository) HasPermission(ctx context.Context, userID uint, permission string) (bool, error) {
	const op = richerror.Op("repository.postgres.user.HasPermission")

	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_permissions
			WHERE user_id = $1 AND permission = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, permission).Scan(&exists)
	if err != nil {
		return false, richerror.New(op).
			WithErr(err).
			WithMessage("failed to check permission")
	}

	return exists, nil
}