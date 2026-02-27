package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) RefreshAccessToken(ctx context.Context, req param.RefreshAccessTokenRequest) (param.RefreshAccessTokenResponse, error) {
	const op = richerror.Op("userservice.RefreshAccessToken")

	logger.Info("Refresh token request received")

	// 1. Validate JWT signature and expiry
	claims, err := s.auth.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		logger.Warn("Invalid refresh token attempt", "error", err.Error())
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("invalid or expired refresh token").
			WithKind(richerror.KindInvalid).
			WithErr(err)
	}

	// 2. Check token is in DB and not revoked
	valid, err := s.repo.IsRefreshTokenValid(ctx, hashToken(req.RefreshToken))
	if err != nil {
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("failed to validate refresh token").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}
	if !valid {
		logger.Warn("Revoked or unknown refresh token used", "user_id", claims.UserID)
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("invalid or expired refresh token").
			WithKind(richerror.KindInvalid)
	}

	// 3. Get user
	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("failed to retrieve user").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	if !user.IsActive {
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("account is deactivated").
			WithKind(richerror.KindForbidden)
	}

	// 4. Revoke old token
	if err := s.repo.RevokeRefreshToken(ctx, hashToken(req.RefreshToken)); err != nil {
		logger.Error("Failed to revoke old refresh token", "user_id", user.ID, "error", err.Error())
	}

	// 5. Issue new access token
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("failed to create access token").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	// 6. Issue new refresh token and save it
	newRefreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("failed to create refresh token").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	if err := s.repo.SaveRefreshToken(
		ctx,
		user.ID,
		hashToken(newRefreshToken),
		time.Now().Add(s.auth.RefreshExpirationTime()),
	); err != nil {
		return param.RefreshAccessTokenResponse{}, richerror.New(op).
			WithMessage("failed to save refresh token").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	logger.Info("Access token refreshed successfully", "user_id", user.ID)

	return param.RefreshAccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}