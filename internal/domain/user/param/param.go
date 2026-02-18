package param

import "time"

type UserInfo struct {
	ID                    uint             `json:"id"`
	Username              string           `json:"username"`
	PhoneNumber           string           `json:"phone_number"`
	AvatarURL             string           `json:"avatar_url"`
	Bio                   string           `json:"bio"`
	Role                  string           `json:"role"`
	SubscriptionTier      string           `json:"subscription_tier"`
	IsActive              bool             `json:"is_active"`
	UsernameLastChangedAt *time.Time       `json:"username_last_changed_at,omitempty"`
	CreatedAt             time.Time        `json:"created_at"`
}

type RegisterRequest struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	User   UserInfo `json:"user"`
	Tokens Tokens
}

type GetProfileResponse struct {
	User UserInfo `json:"user"`
}

type UpdateProfileRequest struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type UpdateProfileResponse struct {
	User UserInfo `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}