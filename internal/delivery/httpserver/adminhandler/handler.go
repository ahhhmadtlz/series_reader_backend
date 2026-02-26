package adminhandler

import (
	"context"

	userparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
)

type UserService interface {
    ChangeUserRole(ctx context.Context, targetUserID uint, req userparam.ChangeUserRoleRequest) (userparam.ChangeUserRoleResponse, error)
    GetUserWithPermissions(ctx context.Context, targetUserID uint) (userparam.AdminUserInfo, error)
    GrantPermission(ctx context.Context, targetUserID uint, req userparam.GrantPermissionRequest, grantedBy uint) (userparam.AdminUserInfo, error)
    RevokePermission(ctx context.Context, targetUserID uint, req userparam.RevokePermissionRequest) (userparam.AdminUserInfo, error)
}

type Handler struct {
    userService UserService
}

func New(userService UserService) Handler {
    return Handler{
        userService: userService,
    }
}