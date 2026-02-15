package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)



func (s Service) ChangePassword(ctx context.Context, userID uint, req param.ChangePasswordRequest) (param.ChangePasswordResponse, error) {
	const op = richerror.Op("userservice.ChangePassword")

	logger.Info("change password request","user_id",userID)

	user,err:=s.repo.GetUserByID(ctx,userID)
	if err !=nil{
		logger.Error("Failed to get user",
	  "user_id",userID,
		"error",err.Error(),
  	)
		return param.ChangePasswordResponse{},richerror.New(op).
		 WithMessage("failed to user").
		 WithKind(richerror.KindUnexpected).
		 WithErr(err)
	}

	if !user.IsActive {
		logger.Warn("Change password failed - account deactivated",
	     "user_id",userID,
	 )
	 return param.ChangePasswordResponse{},richerror.New(op).
	 		WithMessage("account is deactivated").
			WithKind(richerror.KindForbidden)
	}

  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		logger.Warn("Change password failed - incorrect old password",
			"user_id", userID,
		)
		return param.ChangePasswordResponse{}, richerror.New(op).
			WithMessage("old password is incorrect").
			WithKind(richerror.KindInvalid)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash new password",
			"user_id", userID,
			"error", err.Error(),
		)
		return param.ChangePasswordResponse{}, richerror.New(op).
			WithMessage("failed to hash password").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}
	err = s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
	if err != nil {
		logger.Error("Failed to update password",
			"user_id", userID,
			"error", err.Error(),
		)
		return param.ChangePasswordResponse{}, richerror.New(op).
			WithMessage("failed to change password").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	logger.Info("Password changed successfully", "user_id", userID)

	return param.ChangePasswordResponse{
		Message: "password changed successfully",
	}, nil
}