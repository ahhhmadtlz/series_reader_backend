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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}