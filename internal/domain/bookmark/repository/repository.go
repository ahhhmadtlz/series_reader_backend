package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
)

type Repository interface {
	Create(ctx context.Context, bookmark entity.Bookmark)(entity.Bookmark,error)
	GetBookmarksWithSeriesByUserID(ctx context.Context, userID uint) ([]param.BookmarkInfo, error)
	Delete(ctx context.Context,userID uint,seriesID uint)error
}