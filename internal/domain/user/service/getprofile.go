package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetProfile(ctx context.Context, userID uint) (param.GetProfileResponse,error){
	const op = richerror.Op("userservice.GetProfile")

	logger.Info("Get Profile request","user_id",userID)

	user,err:=s.repo.GetUserByID(ctx, userID)

	if err!=nil{
		logger.Error("Failed to get user profile",
	   "user_id",userID,
		 "error",err.Error(),
 	  )
		return param.GetProfileResponse{},richerror.New(op).WithMessage("failed to get profile").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	logger.Info("Profile retrieved successfully","user_id",userID)

	return param.GetProfileResponse{
		User: toUserInfo(user),
	},nil
}
