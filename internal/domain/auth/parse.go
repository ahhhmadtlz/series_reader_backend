package auth

import (
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/golang-jwt/jwt/v5"
)

func (s Service) ParseToken(tokenStr string) (*Claims, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		return nil, richerror.New(OpParseToken).WithMessage("empty or invalid token format").WithKind(richerror.KindInvalid)
	}

	key := []byte(s.config.SignKey)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, richerror.New(OpParseToken).WithMessage("invalid signing method").WithKind(richerror.KindInvalid).WithMeta("algorithm", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, richerror.New(OpParseToken).WithMessage("invalid token").WithKind(richerror.KindInvalid)
}

func (s Service) ParseBearerToken(bearerToken string) (*Claims, error) {
	tokenStr := strings.TrimPrefix(bearerToken, "Bearer ")
	return s.ParseToken(tokenStr)
}

func (s Service) ParseRefreshToken(refreshToken string) (*Claims, error) {
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims.Subject != s.config.RefreshSubject {
		return nil, richerror.New(OpParseToken).WithMessage("invalid refresh token subject").WithKind(richerror.KindInvalid)
	}

	return claims, nil
}