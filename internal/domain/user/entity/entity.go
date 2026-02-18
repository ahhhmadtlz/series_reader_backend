package entity

import (
	"time"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
)

// ============================================
// Permissions - Capabilities
// ============================================
const (
	PermissionModerateComments  = "moderate_comments"
	PermissionUnpublishContent  = "unpublish_content"
	PermissionManageSeriesGlobal = "manage_series_global"
	PermissionManageUsers       = "manage_users"
	PermissionViewAnalytics     = "view_analytics"
)


// ============================================
// User Entity
// ============================================

type User struct {
	ID                    uint
	Username              string
	PhoneNumber           string
	Password              string
	AvatarURL             string
	Bio                   string
	Role                  sharedentity.Role
	SubscriptionTier      sharedentity.SubscriptionTier
	IsActive              bool
	UsernameLastChangedAt *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (u User) IsAdmin() bool {
	return u.Role == sharedentity.AdminRole
}

func (u User) IsManager() bool {
	return u.Role == sharedentity.ManagerRole
}

func (u User) IsPremium() bool {
	return u.SubscriptionTier == sharedentity.PremiumTier
}