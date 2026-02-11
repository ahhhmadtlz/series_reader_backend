package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) AddPages(ctx context.Context, req param.AddChapterPagesRequest)error {
	const op =richerror.Op("service.chapter.AddPages")

	//Verify chapter exists

	_,err:=s.repo.GetByID(ctx,req.ChapterID)

	if err !=nil{
		return richerror.New(op).WithErr(err).WithMessage("chapter not found").WithKind(richerror.KindNotFound)
	}

	//conver req to entities

	pages:=make([]entity.ChapterPage,len(req.Pages))

	for i,p :=range req.Pages{
		pages[i]=entity.ChapterPage{
			ChapterID: req.ChapterID,
			PageNumber: p.PageNumber,
			ImageURL: p.ImageURL,
		}
	}

	//save pages to database
	err =s.repo.CreatePages(ctx,pages)
	if err !=nil{
		return richerror.New(op).WithErr(err).WithMessage("failed to add pages").WithKind(richerror.KindUnexpected)
	}
	return nil
}

func (s Service) GetPages(ctx context.Context, chapterID uint) ([]param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.GetPages")

	pages, err := s.repo.GetPagesByChapterID(ctx, chapterID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get pages").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.ChapterPageResponse, len(pages))
	for i, p := range pages {
		responses[i] = param.ChapterPageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
		}
	}

	return responses, nil
}

func (s Service) GetPageByNumber(ctx context.Context, chapterID uint, pageNumber int) (param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.GetPageByNumber")

	page, err := s.repo.GetPageByNumber(ctx, chapterID, pageNumber)
	if err != nil {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("page not found").
			WithKind(richerror.KindNotFound)
	}

	return param.ChapterPageResponse{
		ID:         page.ID,
		PageNumber: page.PageNumber,
		ImageURL:   page.ImageURL,
	}, nil
}