package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) UpdateProfile(ctx context.Context,userID uint, req param.UpdateProfileRequest)(param.UpdateProfileResponse,error){
	const op = richerror.Op("userservice.UpdateProfile")

	logger.Info("Update profile request",
    "user_id",userID,
		"username",req.Username,
  )
  
	currentUser, err :=s.repo.GetUserByID(ctx,userID)

	if err !=nil{
		logger.Error("Failed to get user",
	   "user_id",userID,
		 "error",err.Error(),
  	)
		return param.UpdateProfileResponse{},richerror.New(op).
		WithMessage("failed to get user").
		WithKind(richerror.KindUnexpected).
		WithErr(err)
	}

	// Check if user is active
	if !currentUser.IsActive {
		logger.Warn("Update profile failed - account deactivated",
			"user_id", userID,
		)
		return param.UpdateProfileResponse{}, richerror.New(op).
			WithMessage("account is deactivated").
			WithKind(richerror.KindForbidden)
	}

	// Prepare updated user
	updatedUser := entity.User{
		ID:                    userID,
		Username:              req.Username,
		AvatarURL:             req.AvatarURL,
		Bio:                   req.Bio,
		UsernameLastChangedAt: currentUser.UsernameLastChangedAt,
	}


	// If username changed, update the timestamp
	if req.Username != currentUser.Username {
		now := time.Now()
		updatedUser.UsernameLastChangedAt = &now
		logger.Info("Username changed",
			"user_id", userID,
			"old_username", currentUser.Username,
			"new_username", req.Username,
		)
	}

		// Update in database
	user, err := s.repo.UpdateUser(ctx, updatedUser)
	if err != nil {
		logger.Error("Failed to update profile",
			"user_id", userID,
			"error", err.Error(),
		)
		return param.UpdateProfileResponse{}, richerror.New(op).
			WithMessage("failed to update profile").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

		logger.Info("Profile updated successfully",
		"user_id", userID,
		"username", user.Username,
	)

	return param.UpdateProfileResponse{
		User: toUserInfo(user),
	}, nil

}