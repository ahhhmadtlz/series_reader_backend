package auth

import (
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/golang-jwt/jwt/v5"
)

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user, s.config.RefreshSubject, s.config.RefreshExpirationTime)
}

func (s Service) createToken(
	user entity.User,
	subject string,
	expireDuration time.Duration,
) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
		Role:user.Role,
		SubscriptionTier: user.SubscriptionTier,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}