package param

import "time"

type UserInfo struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url"`
	Bio         string    `json:"bio"`
	IsActive    bool `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}


type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User UserInfo `json:"user"`
		Tokens Tokens
}

 