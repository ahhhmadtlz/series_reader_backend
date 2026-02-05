package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
)

type Repository interface {
	Create(ctx context.Context, series entity.Series)(entity.Series,error)
	GetByID(ctx context.Context,id uint)(entity.Series,error)
	GetBySlug(ctx context.Context,slug string)(entity.Series,error)
	GetList(ctx context.Context,req param.GetListRequest)([]entity.Series,int,error)
	Update(ctx context.Context,id uint,series entity.Series)(entity.Series,error)
	Delete(ctx context.Context,id uint)error
	IncrementViewCount(ctx context.Context,id uint)error
	IsSlugExists(ctx context.Context,slug string)(bool,error)
}
