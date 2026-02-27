package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) Logout(ctx context.Context, req param.LogoutRequest) error {
	const op = richerror.Op("userservice.Logout")

	if err := s.repo.RevokeRefreshToken(ctx, hashToken(req.RefreshToken)); err != nil {
		logger.Error("Failed to revoke refresh token on logout", "error", err.Error())
		return richerror.New(op).
			WithMessage("failed to logout").
			WithKind(richerror.KindUnexpected).
			WithErr(err)
	}

	logger.Info("User logged out", "user_id", req.UserID)
	return nil
}