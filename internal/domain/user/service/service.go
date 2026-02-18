package service

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	userRepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/repository"
)

type AuthService interface {
	CreateAccessToken(user entity.User)(string,error)
	CreateRefreshToken(user entity.User)(string,error)
	ParseRefreshToken(refreshToken string)(*auth.Claims,error)
}



type Service struct {
	repo userRepository.Repository
	auth AuthService
}

func New(authService AuthService,repo userRepository.Repository)Service{
	return Service{
		auth:authService,
		repo:repo,
	}
}

func toUserInfo(user entity.User) param.UserInfo {
	return param.UserInfo{
		ID:                    user.ID,
		Username:              user.Username,
		PhoneNumber:           user.PhoneNumber,
		AvatarURL:             user.AvatarURL,
		Bio:                   user.Bio,
		Role:                  user.Role.String(),
		SubscriptionTier:      user.SubscriptionTier.String(),
		IsActive:              user.IsActive,
		UsernameLastChangedAt: user.UsernameLastChangedAt,
		CreatedAt:             user.CreatedAt,
	}
}


func toAdminUserInfo(user entity.User, permissions []string) param.AdminUserInfo {
	return param.AdminUserInfo{
		ID:               user.ID,
		Username:         user.Username,
		PhoneNumber:      user.PhoneNumber,
		Role:             user.Role.String(),
		SubscriptionTier: user.SubscriptionTier.String(),
		IsActive:         user.IsActive,
		Permissions:      permissions,
		CreatedAt:        user.CreatedAt,
	}
}