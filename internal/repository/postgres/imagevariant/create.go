package imagevariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) CreateVariant(ctx context.Context, variant entity.ImageVariant) (entity.ImageVariant, error) {
	const op = richerror.Op("postgres.imagevariant.CreateVariant")

	const query = `
		INSERT INTO image_variants (chapter_page_id, kind, image_url, remote_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	row :=r.db.QueryRowContext(ctx,query,
	   variant.ChapterPageID,
		 variant.Kind,
		 variant.ImageURL,
		 variant.RemotePath,
	)

	if err:=row.Scan(&variant.ID,&variant.CreatedAt);err!=nil{
		logger.Error("failed to insert image variant","error",err)
		return entity.ImageVariant{},richerror.New(op).
		   WithErr(err).
			 WithMessage("failed to create image variant").
			 WithKind(richerror.KindUnexpected)
	}

	logger.Debug("image variant created", "id", variant.ID, "kind", variant.Kind)

	return variant, nil
	

}