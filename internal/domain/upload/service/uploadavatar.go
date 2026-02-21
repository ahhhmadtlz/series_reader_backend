package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)


func (s Service) UploadAvatar(ctx context.Context, req param.UploadAvatarRequest) (param.UploadAvatarResponse, error) {
	const op = richerror.Op("service.upload.UploadAvatar")

	logger.Info("Upload avatar request",
		"user_id", req.UserID,
		"filename", req.Header.Filename,
		"size", req.Header.Size,
	)

	user, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		logger.Error("Failed to get user",
			"user_id", req.UserID,
			"error", err.Error(),
		)
		return param.UploadAvatarResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get user").
			WithKind(richerror.KindUnexpected)
	}

	if !user.IsActive {
		logger.Warn("Upload avatar failed - account deactivated",
			"user_id", req.UserID,
		)
		return param.UploadAvatarResponse{}, richerror.New(op).
			WithMessage("account is deactivated").
			WithKind(richerror.KindForbidden)
	}

	// 2. Delete old avatar if exists (best effort)
	oldAvatar, err := s.uploadRepo.GetLatestByOwner(ctx, req.UserID, entity.ImageKindAvatar)
	if err == nil {
		// Old avatar exists, delete it
		logger.Info("Deleting old avatar",
			"user_id", req.UserID,
			"old_avatar_id", oldAvatar.ID,
			"old_path", oldAvatar.StoredPath,
		)

		// Delete from DB first
		_ = s.uploadRepo.DeleteByID(ctx, oldAvatar.ID)

		// Delete from storage (best effort, log if fails)
		if err := s.storage.Delete(ctx, oldAvatar.StoredPath); err != nil {
			logger.Error("Failed to delete old avatar file",
				"user_id", req.UserID,
				"stored_path", oldAvatar.StoredPath,
				"error", err.Error(),
			)
			// Don't fail the upload if old file deletion fails
		}
	}

	saveReq := storage.SaveRequest{
		File:     req.File,
		Filename: req.Header.Filename,
		OwnerID:  req.UserID,
		Kind:     entity.ImageKindAvatar,
		MimeType: req.Header.Header.Get("Content-Type"),
		Size:     req.Header.Size,
	}

	result, err := s.storage.Save(ctx, saveReq)
	if err != nil {
		logger.Error("Failed to save avatar to storage",
			"user_id", req.UserID,
			"error", err.Error(),
		)
		return param.UploadAvatarResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save avatar file").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Avatar saved to storage",
		"user_id", req.UserID,
		"stored_path", result.StoredPath,
		"url", result.URL,
	)

	// 4. Save avatar record to DB
	avatarImg := entity.UploadedImage{
		OwnerID:    req.UserID,
		Kind:       entity.ImageKindAvatar,
		Filename:   req.Header.Filename,
		StoredPath: result.StoredPath,
		URL:        result.URL,
		MimeType:   saveReq.MimeType,
		SizeBytes:  req.Header.Size,
	}

	savedImg, err := s.uploadRepo.Save(ctx, avatarImg)
	if err != nil {
		logger.Error("Failed to save avatar to database",
			"user_id", req.UserID,
			"error", err.Error(),
		)

		// CRITICAL: Rollback - delete file from storage
		_ = s.storage.Delete(ctx, result.StoredPath)

		return param.UploadAvatarResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save avatar record").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Avatar record saved to database",
		"user_id", req.UserID,
		"image_id", savedImg.ID,
	)

	// 5. Update user's avatar_url
	user.AvatarURL = savedImg.URL
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		logger.Error("Failed to update user avatar URL",
			"user_id", req.UserID,
			"error", err.Error(),
		)

		// CRITICAL: Rollback - delete both DB record and file
		_ = s.uploadRepo.DeleteByID(ctx, savedImg.ID)
		_ = s.storage.Delete(ctx, result.StoredPath)

		return param.UploadAvatarResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to update user profile").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Avatar uploaded successfully",
		"user_id", req.UserID,
		"avatar_url", savedImg.URL,
	)

	return param.UploadAvatarResponse{
		AvatarURL: savedImg.URL,
	}, nil

}