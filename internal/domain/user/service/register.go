package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {
	const op = richerror.Op("userservice.Register")

	logger.Info("User registration attempt",
		"phone_number", req.PhoneNumber,
		"username", req.Username,
	)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "error", err.Error())
		return param.RegisterResponse{}, richerror.New(op).
			WithMessage("failed to hash password").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	user := entity.User{
		Username:    req.Username,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
		IsActive:    true,
	}

	createdUser, err := s.repo.RegisterUser(ctx, user)
	if err != nil {
		logger.Error("Failed to create user",
			"phone_number", req.PhoneNumber,
			"error", err.Error(),
		)
		return param.RegisterResponse{}, richerror.New(op).
			WithMessage("failed to register user").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	logger.Info("User registered successfully",
		"user_id", createdUser.ID,
		"phone_number", createdUser.PhoneNumber,
	)

	return param.RegisterResponse{
		User: toUserInfo(createdUser),
	}, nil
}