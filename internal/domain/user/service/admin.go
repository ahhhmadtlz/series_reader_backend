package service

import (
	"context"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) ChangeUserRole(ctx context.Context, targetUserID uint, req param.ChangeUserRoleRequest) (param.ChangeUserRoleResponse, error) {
	const op = richerror.Op("service.user.ChangeUserRole")

	role, ok := sharedentity.MapToRoleEntity(req.Role)
	if !ok {
		return param.ChangeUserRoleResponse{}, richerror.New(op).
			WithMessage("invalid role").
			WithKind(richerror.KindInvalid)
	}

	err := s.repo.UpdateUserRole(ctx, targetUserID, role)
	if err != nil {
		return param.ChangeUserRoleResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to change user role")
	}

	return param.ChangeUserRoleResponse{
		Message: "user role updated successfully",
		UserID:  targetUserID,
		Role:    role.String(),
	}, nil
}

func (s Service) GrantPermission(ctx context.Context, targetUserID uint, req param.GrantPermissionRequest, grantedBy uint) (param.AdminUserInfo, error) {
	const op = richerror.Op("service.user.GrantPermission")

	err := s.repo.GrantPermission(ctx, targetUserID, req.Permission, grantedBy)
	if err != nil {
		return param.AdminUserInfo{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to grant permission")
	}

	return s.getAdminUserInfo(ctx, op, targetUserID)
}

func (s Service) RevokePermission(ctx context.Context, targetUserID uint, req param.RevokePermissionRequest) (param.AdminUserInfo, error) {
	const op = richerror.Op("service.user.RevokePermission")

	err := s.repo.RevokePermission(ctx, targetUserID, req.Permission)
	if err != nil {
		return param.AdminUserInfo{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to revoke permission")
	}

	return s.getAdminUserInfo(ctx, op, targetUserID)
}

func (s Service) GetUserWithPermissions(ctx context.Context, targetUserID uint) (param.AdminUserInfo, error) {
	const op = richerror.Op("service.user.GetUserWithPermissions")
	return s.getAdminUserInfo(ctx, op, targetUserID)
}

func (s Service) getAdminUserInfo(ctx context.Context, op richerror.Op, userID uint) (param.AdminUserInfo, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return param.AdminUserInfo{}, richerror.New(op).
			WithErr(err).
			WithMessage("user not found").
			WithKind(richerror.KindNotFound)
	}

	permissions, err := s.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return param.AdminUserInfo{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get permissions")
	}

	return toAdminUserInfo(user, permissions), nil
}
