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
	CreatePages(ctx context.Context, pages []entity.ChapterPage) ([]entity.ChapterPage, error)
	GetPagesByChapterID(ctx context.Context, chapterID uint)([]entity.ChapterPage,error)
	GetPageByNumber(ctx context.Context,chapterID uint ,pageNumber int)(*entity.ChapterPage,error)
	GetPageByID(ctx context.Context, pageID uint) (*entity.ChapterPage, error)
	DeletePage(ctx context.Context,pageID uint)error
	UpdatePageNumbers(ctx context.Context, updates []entity.PageNumberUpdate) error 

	GetChapterWithPages(ctx context.Context,chapterID uint)(*entity.Chapter,[]entity.ChapterPage,error)

}