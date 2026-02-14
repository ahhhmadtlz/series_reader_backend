package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
)

type Repository interface {
	RegisterUsr(ctx context.Context,user entity.User)(entity.User,error)
	GetUserByID(ctx context.Context,userID uint)(entity.User,error)
	GetUserByPhoneNumber(ctx context.Context,phoneNumber string)(entity.User,error)

	IsPhoneNumberUnique(ctx context.Context, phoneNumber string)(bool,error)
	IsUsernameUnique(ctx context.Context,username string)(bool,error)
}