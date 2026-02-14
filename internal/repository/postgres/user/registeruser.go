package user

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) Registeruser(ctx context.Context,user entity.User)(entity.User,error){
	const op =richerror.Op("repository.postgres.user.RegisterUser")

	query:=`
	  INSERT INTO users(
		 username, phone_number, password, avatar_url, bio, is_active
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var createdUser entity.User

	err :=r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.PhoneNumber,
		user.Password,
		user.AvatarURL,
		user.Bio,
		true,
	).Scan(&createdUser.ID,&createdUser.CreatedAt,&createdUser.UpdatedAt)

	if err !=nil{
		return entity.User{},richerror.New(op).WithErr(err).WithMessage("failed to create user")
	}

	createdUser.Username=user.Username
	createdUser.PhoneNumber=user.PhoneNumber
	createdUser.Password=user.Password
	createdUser.AvatarURL=user.AvatarURL
	createdUser.Bio=user.Bio
	createdUser.IsActive=true

	return createdUser,nil
}