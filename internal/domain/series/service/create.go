package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/slugify"
	"github.com/segmentio/ksuid"
)

func (s Service) Create(ctx context.Context,req param.CreateSeriesRequest)(param.SeriesResponse,error){
	const op=richerror.Op("service.series.Create")

	slug:=req.Slug
	if slug==""{
		slug=slugify.Make(req.Title)
	}else{
		slug=slugify.Make(slug)
	}

	exists,err:=s.repo.IsSlugExists(ctx,slug)

	if err!=nil{
		return param.SeriesResponse{},richerror.New(op).WithErr(err).WithMessage("failedt o check slug existence").WithKind(richerror.KindUnexpected)
	}

	if exists {
		slug =slugify.MakeUnique(req.Title,func(testSlug string)bool{
			exists,_:=s.repo.IsSlugExists(ctx,testSlug)
			return exists
		})
	}
	slugID := ksuid.New().String()[0:8]

	fullSlug := slug + "-" + slugID

	series := entity.Series{
		Title:             req.Title,
		Slug:              slug,
		SlugID:            slugID,
		FullSlug:          fullSlug,
		Description:       req.Description,
		Author:            req.Author,
		Artist:            req.Artist,
		Status:            req.Status,
		Type:              req.Type,
		Genres:            req.Genres,
		AlternativeTitles: req.AlternativeTitles,
		CoverImageURL:     req.CoverImageURL,
		PublicationYear:   req.PublicationYear,
		ViewCount:        	0,
		Rating:             0.0,
		IsPublished:       req.IsPublished,
	}

	// Ensure empty slices are initialized
	if series.Genres == nil {
		series.Genres = []string{}
	}
	if series.AlternativeTitles == nil {
		series.AlternativeTitles = []string{}
	}

	createdSeries, err := s.repo.Create(ctx, series)
	if err != nil {
		return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to create series").WithKind(richerror.KindUnexpected)
	}

	return toSeriesResponse(createdSeries), nil
}


