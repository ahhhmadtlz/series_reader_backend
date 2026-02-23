package service

import (
	"context"
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) BulkUploadPages(ctx context.Context,req param.BulkUploadParam) ([]param.ChapterPageResponse, error) {
  const op = richerror.Op("service.chapter.BulkUploadPages")

	// 1. verify chapter exists
	_,err:=s.repo.GetByID(ctx,req.ChapterID)
	if err !=nil{
		return nil,richerror.New(op).
		 WithErr(err).
		 WithMessage("chapter not found").
		 WithKind(richerror.KindNotFound)
	}
	if len(req.Files)==0 {
		return nil,richerror.New(op).
		 WithMessage("no files provided").
		 WithKind(richerror.KindInvalid)
	} 

	//2. get existing pages to determine next page number
	existing, err :=s.repo.GetPagesByChapterID(ctx,req.ChapterID)
	if err !=nil{
		return nil,richerror.New(op).
		 WithErr(err).
		 WithMessage("failed to get existing pages").
		 WithKind(richerror.KindUnexpected)
	}
	nextPageNumber :=len(existing)+1

	//3. upload each file and collect pages
	var pages []entity.ChapterPage
	var storedPaths []string //for rollback

	for i , fileHeader :=range req.Files{
		file,err :=fileHeader.Open()
		if err !=nil{
			//rollback all saved files
			for _,path :=range storedPaths {
				_=s.storage.Delete(ctx,path)
			}
			return nil,richerror.New(op).
			  WithErr(err).
				WithMessage(fmt.Sprintf("failed to open file %d",i+1)).
				WithKind(richerror.KindUnexpected)
		}
		result,err :=s.storage.Save(ctx,storage.SaveRequest{
			File: file,
			Filename: fileHeader.Filename,
			OwnerID: req.ChapterID,
		  Kind:     uploadentity.ImageKindChapterPage,
			MimeType: fileHeader.Header.Get("Content-Type"),
			Size: fileHeader.Size,
		})
		file.Close()
		if err!=nil{
			//rollback all saved files
			for _,path:=range storedPaths {
				_=s.storage.Delete(ctx,path)
			}
		return nil, richerror.New(op).
				WithErr(err).
				WithMessage(fmt.Sprintf("failed to save file %d", i+1)).
				WithKind(richerror.KindUnexpected)
		}
		storedPaths=append(storedPaths, result.StoredPath)
		pages=append(pages, entity.ChapterPage{
			ChapterID: req.ChapterID,
			PageNumber: nextPageNumber + i,
			ImageURL: result.URL,
			RemotePath: result.StoredPath,
		})
	}
	//4. save all records in one DB call (uses transaction in repo)

created, err := s.repo.CreatePages(ctx, pages)
if err != nil {
    for _, path := range storedPaths {
        _ = s.storage.Delete(ctx, path)
    }
    return nil, richerror.New(op).
        WithMessage("failed to save pages to database").
        WithKind(richerror.KindUnexpected)
}

	logger.Info("bulk pages uploaded","chapter_id",req.ChapterID,"count",len(pages))

responses := make([]param.ChapterPageResponse, len(created))
for i, p := range created {
    responses[i] = param.ChapterPageResponse{
        ID:         p.ID,
        PageNumber: p.PageNumber,
        ImageURL:   p.ImageURL,
    }
}

	return responses,nil

}
