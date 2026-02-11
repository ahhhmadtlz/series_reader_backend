package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
)

type Repository interface {
	//chapter CRUD
	Create(ctx context.Context,chapter *entity.Chapter)(*entity.Chapter,error)
	GetByID(ctx context.Context , id uint)(*entity.Chapter,error)
	GetBySeriesID(ctx context.Context,seriesID uint)([]*entity.Chapter,error)
	// Update(ctx context.Context,id uint,chapter *entity.Chapter)(*entity.Chapter,error)
	Delete(ctx context.Context,id uint)error

	//Page operation
	CreatePages(ctx context.Context,pages []entity.ChapterPage)error
	GetPagesByChapterID(ctx context.Context, chapterID uint)([]entity.ChapterPage,error)
	GetPageByNumber(ctx context.Context,chapterID uint ,pageNumber int)(*entity.ChapterPage,error)
	DeletePage(ctx context.Context,pageID uint)error

	//Optimized combined query
	GetChapterWithPages(ctx context.Context,chapterID uint)(*entity.Chapter,[]entity.ChapterPage,error)

}