package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(ctx context.Context, req param.LoginRequest)(param.LoginResponse,error){
	const op=richerror.Op("userservice.Login")

	logger.Info("Login attempt","phone_number",req.PhoneNumber)

	user, err:=s.repo.GetUserByPhoneNumber(ctx,req.PhoneNumber)

	if err !=nil{
		logger.Warn("Login failed - user not found", "phone_number",req.PhoneNumber)
		return param.LoginResponse{},richerror.New(op).
		WithMessage("phone number or password is incorrect").
		WithKind(richerror.KindInvalid)
	}


	if !user.IsActive {
			logger.Warn("Login failed - account deactivated",
			"user_id",user.ID,
			"phone_number",user.PhoneNumber,
	  )

		return param.LoginResponse{},richerror.New(op).
			WithMessage("account is deactivated").
			WithKind(richerror.KindForbidden)
	}


	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password))

	if err !=nil{
		  logger.Warn("Login failed - incorrect password",
	    "user_id",user.ID,
	  	"phone_number",user.PhoneNumber,
  	)

		return param.LoginResponse{},richerror.New(op).
	  	WithMessage("phone number or password is incorrect").
		  WithKind(richerror.KindInvalid)
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		logger.Error("Failed to create access token",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return param.LoginResponse{}, richerror.New(op).
			WithMessage("failed to create access token").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		logger.Error("Failed to create refresh token",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return param.LoginResponse{}, richerror.New(op).
			WithMessage("failed to create refresh token").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	logger.Info("User logged in successfully",
		"user_id", user.ID,
		"phone_number", user.PhoneNumber,
	)

	return param.LoginResponse{
		User:toUserInfo(user),
		Tokens:param.Tokens{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		},
	},nil
}