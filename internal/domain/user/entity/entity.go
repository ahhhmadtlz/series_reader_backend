package entity

import "time"

type User struct {
	ID          uint
	Username    string
	PhoneNumber string
	Password    string
	AvatarURL   string
	Bio         string
	IsActive    bool
	UsernameLastChangedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}