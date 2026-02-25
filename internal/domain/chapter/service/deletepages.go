package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) DeletePage(ctx context.Context, chapterID uint, pageNumber int) error {
	const op = richerror.Op("service.chapter.DeletePage")

	// 1. Get the page to retrieve its remote path
	page, err := s.repo.GetPageByNumber(ctx, chapterID, pageNumber)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("page not found").
			WithKind(richerror.KindNotFound)
	}

	//2. fetch variants before deleting the page
  variants,err:=s.variantRepo.GetVariantsByPageID(ctx,page.ID)
	if err !=nil{
		logger.Error("failed to get variants for page","page_id",page.ID,"error",err)
	}

 //3. delete variant files from disk
 for _, v:= range variants {
	 if v.RemotePath !=""{
		if err :=s.storage.Delete(ctx, v.RemotePath);err !=nil{
			logger.Error("failed to delete variant file","rempte_path",v.RemotePath,"error",err)
		}
	 }
 }

 //4. delete variant rows from DB
 if err :=s.variantRepo.DeleteVariantsByPageID(ctx,page.ID); err !=nil{
	logger.Error("failed to delete variant rows","page_id",page.ID,"error",err)
 }

	// 6. delete page row from DB
	if err := s.repo.DeletePage(ctx, page.ID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete page").
			WithKind(richerror.KindUnexpected)
	}

	// 3. Delete file from storage (best effort — log but don't fail)
	if page.RemotePath != "" {
		if err := s.storage.Delete(ctx, page.RemotePath); err != nil {
			logger.Error("failed to delete page file from storage",
				"chapter_id", chapterID,
				"page_number", pageNumber,
				"remote_path", page.RemotePath,
				"error", err,
			)
		}
	}

	logger.Info("page deleted", "chapter_id", chapterID, "page_number", pageNumber)

	return nil
}