package upload

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) Save(ctx context.Context,img entity.UploadedImage)(entity.UploadedImage,error){
	const op=richerror.Op("repository.postgres.upload.save")

	query:=`
		INSERT INTO uploaded_images (
			owner_id,
			kind,
			filename,
			stored_path,
			url,
			mime_type,
			size_bytes
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id,owner_id,kind,filename,stored_path,url,mime_type,size_bytes, created_at
	`

	var savedImg entity.UploadedImage

	err :=r.db.QueryRowContext(
      ctx,
			query,
			img.OwnerID,
			img.Kind,
			img.Filename,
			img.StoredPath,
			img.URL,
			img.MimeType,
			img.SizeBytes,
		).Scan(
			&savedImg.ID,
			&savedImg.OwnerID,
			&savedImg.Kind,
			&savedImg.Filename,
			&savedImg.StoredPath,
			&savedImg.URL,
			&savedImg.MimeType,
			&savedImg.SizeBytes,
			&savedImg.CreatedAt,
		)

		if err !=nil{
			return  entity.UploadedImage{},richerror.New(op).
			        WithErr(err).
							WithMessage("failed to save uploaded image").
							WithKind(richerror.KindUnexpected)
		}

		return savedImg, nil

}
