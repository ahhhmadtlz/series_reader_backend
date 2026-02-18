package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
)

type Repository interface {
	RegisterUser(ctx context.Context,user entity.User)(entity.User,error)
	GetUserByID(ctx context.Context,userID uint)(entity.User,error)
	GetUserByPhoneNumber(ctx context.Context,phoneNumber string)(entity.User,error)

	IsPhoneNumberUnique(ctx context.Context, phoneNumber string)(bool,error)
	IsUsernameUnique(ctx context.Context,username string)(bool,error)

	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)  
	UpdatePassword(ctx context.Context, userID uint, newPassword string) error

	GetUserPermissions(ctx context.Context, userID uint) ([]string, error)
	HasPermission(ctx context.Context, userID uint, permission string) (bool, error) 

	//admin
	UpdateUserRole(ctx context.Context,userID uint ,role sharedentity.Role)error
	GrantPermission(ctx context.Context,userID uint, permission string, grantedBy uint)error
	RevokePermission(ctx context.Context,userID uint, permission string)error
}