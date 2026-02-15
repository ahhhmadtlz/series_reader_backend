package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
)

type AuthService interface {
	CreateAccessToken(user entity.User)(string,error)
	CreateRefreshToken(user entity.User)(string,error)
	ParseRefreshToken(refreshToken string)(*auth.Claims,error)
}

type Repository interface {
	RegisterUser(ctx context.Context,user entity.User)(entity.User,error)
	GetUserByID(ctx context.Context,userID uint)(entity.User,error)
	GetUserByPhoneNumber(ctx context.Context,phoneNumber string)(entity.User,error)
}

type Service struct {
	repo Repository
	auth AuthService
}

func New(authService AuthService,repo Repository)Service{
	return Service{
		auth:authService,
		repo:repo,
	}
}

func toUserInfo(user entity.User) param.UserInfo {
	return param.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		AvatarURL:   user.AvatarURL,
		Bio:         user.Bio,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
	}
}