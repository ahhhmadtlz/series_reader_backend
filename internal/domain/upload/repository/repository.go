package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
)

type Repository interface {
	Save(ctx context.Context,img entity.UploadedImage)(entity.UploadedImage,error)
	GetByID(ctx context.Context,id uint)(entity.UploadedImage,error)
	GetByOwner(ctx context.Context,ownerID uint, kind entity.ImageKind)([]entity.UploadedImage,error)
  GetLatestByOwner(ctx context.Context, ownerID uint, kind entity.ImageKind) (entity.UploadedImage, error)
	DeleteByID(ctx context.Context,id uint) error
	DeleteByStoredPath(ctx context.Context,storedPath string)error
}

