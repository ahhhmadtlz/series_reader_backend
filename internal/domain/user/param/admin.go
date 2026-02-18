package param

import "time"

// AdminUserInfo - what admin sees about a user
type AdminUserInfo struct {
	ID               uint      `json:"id"`
	Username         string    `json:"username"`
	PhoneNumber      string    `json:"phone_number"`
	Role             string    `json:"role"`
	SubscriptionTier string    `json:"subscription_tier"`
	IsActive         bool      `json:"is_active"`
	Permissions      []string  `json:"permissions"`
	CreatedAt        time.Time `json:"created_at"`
}

type ChangeUserRoleRequest struct {
	Role string `json:"role"`
}

type ChangeUserRoleResponse struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
	Role    string `json:"role"`
}

type GrantPermissionRequest struct {
	Permission string `json:"permission"`
}

type RevokePermissionRequest struct {
	Permission string `json:"permission"`
}