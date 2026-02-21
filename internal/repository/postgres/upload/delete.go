package upload

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) DeleteByID(ctx context.Context, id uint) error {
	const op = richerror.Op("repository.postgres.upload.DeleteByID")

	query := `DELETE FROM uploaded_images WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete uploaded image").
			WithKind(richerror.KindUnexpected)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected").
			WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("uploaded image not found").
			WithKind(richerror.KindNotFound)
	}

	return nil
}

func (r *PostgresRepository) DeleteByStoredPath(ctx context.Context, storedPath string) error {
	const op = richerror.Op("repository.postgres.upload.DeleteByStoredPath")

	query := `DELETE FROM uploaded_images WHERE stored_path = $1`

	result, err := r.db.ExecContext(ctx, query, storedPath)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete uploaded image by path").
			WithKind(richerror.KindUnexpected)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected").
			WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("uploaded image not found by path").
			WithKind(richerror.KindNotFound)
	}

	return nil
}


