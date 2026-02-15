package user

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) UpdateUser(ctx context.Context,user entity.User)(entity.User,error){
	const op=richerror.Op("repository.postgres.user.UpdateUser")

	query:=`
	UPDATE users SET
	  username = $1,
		avatar_url = $2,
		bio = $3,
		username_last_changed_at = $4
	WHERE id = $5
  RETURNING id, username, phone_number, password, avatar_url, bio, is_active, username_last_changed_at, created_at, updated_at
	`
	var updatedUser entity.User

	err :=r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.AvatarURL,
		user.Bio,
		user.UsernameLastChangedAt,
	  user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.PhoneNumber,
		&updatedUser.Password,
		&updatedUser.AvatarURL,
		&updatedUser.Bio,
		&updatedUser.IsActive,
		&updatedUser.UsernameLastChangedAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err == sql.ErrNoRows{
		return  entity.User{},richerror.New(op).WithErr(err).WithMessage("user not found").WithKind(richerror.KindNotFound)
	}

	if err !=nil{
		return  entity.User{},richerror.New(op).WithErr(err).WithMessage("failed to  update user")
	}
	return  updatedUser,nil
}

func (r *PostgresRepository) UpdatePassword(ctx context.Context, userID uint, newPassword string) error {
	const op = richerror.Op("repository.postgres.user.UpdatePassword")

	query := `UPDATE users SET password = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, newPassword, userID)
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to update password")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).WithMessage("user not found").WithKind(richerror.KindNotFound)
	}

	return nil
}