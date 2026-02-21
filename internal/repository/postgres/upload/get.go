package upload

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetByID(ctx context.Context, id uint) (entity.UploadedImage, error) {
	const op = richerror.Op("repository.postgres.upload.GetByID")

	query := `
		SELECT id, owner_id, kind, filename, stored_path, url, mime_type, size_bytes, created_at
		FROM uploaded_images
		WHERE id = $1
	`

	var img entity.UploadedImage

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&img.ID,
		&img.OwnerID,
		&img.Kind,
		&img.Filename,
		&img.StoredPath,
		&img.URL,
		&img.MimeType,
		&img.SizeBytes,
		&img.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return entity.UploadedImage{}, richerror.New(op).
			WithErr(err).
			WithMessage("uploaded image not found").
			WithKind(richerror.KindNotFound)
	}

	if err != nil {
		return entity.UploadedImage{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get uploaded image").
			WithKind(richerror.KindUnexpected)
	}

	return img, nil
}

func (r *PostgresRepository) GetByOwner(ctx context.Context, ownerID uint, kind entity.ImageKind) ([]entity.UploadedImage, error) {
	const op = richerror.Op("repository.postgres.upload.GetByOwner")

	query := `
		SELECT id, owner_id, kind, filename, stored_path, url, mime_type, size_bytes, created_at
		FROM uploaded_images
		WHERE owner_id = $1 AND kind = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID, kind)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to query uploaded images").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var images []entity.UploadedImage

	for rows.Next() {
		var img entity.UploadedImage
		err := rows.Scan(
			&img.ID,
			&img.OwnerID,
			&img.Kind,
			&img.Filename,
			&img.StoredPath,
			&img.URL,
			&img.MimeType,
			&img.SizeBytes,
			&img.CreatedAt,
		)
		if err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan uploaded image").
				WithKind(richerror.KindUnexpected)
		}
		images = append(images, img)
	}

	if err = rows.Err(); err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("error iterating uploaded images").
			WithKind(richerror.KindUnexpected)
	}

	return images, nil
}

func (r *PostgresRepository) GetLatestByOwner(ctx context.Context, ownerID uint, kind entity.ImageKind) (entity.UploadedImage, error) {
	const op = richerror.Op("repository.postgres.upload.GetLatestByOwner")

	query := `
		SELECT id, owner_id, kind, filename, stored_path, url, mime_type, size_bytes, created_at
		FROM uploaded_images
		WHERE owner_id = $1 AND kind = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var img entity.UploadedImage

	err := r.db.QueryRowContext(ctx, query, ownerID, kind).Scan(
		&img.ID,
		&img.OwnerID,
		&img.Kind,
		&img.Filename,
		&img.StoredPath,
		&img.URL,
		&img.MimeType,
		&img.SizeBytes,
		&img.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return entity.UploadedImage{}, richerror.New(op).
			WithErr(err).
			WithMessage("no uploaded image found for owner").
			WithKind(richerror.KindNotFound)
	}

	if err != nil {
		return entity.UploadedImage{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get latest uploaded image").
			WithKind(richerror.KindUnexpected)
	}

	return img, nil
}