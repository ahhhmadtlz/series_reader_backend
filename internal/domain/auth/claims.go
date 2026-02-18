package auth


import (
	"github.com/golang-jwt/jwt/v5"
	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID           uint                        `json:"user_id"`
	Role             sharedentity.Role           `json:"role"`
	SubscriptionTier sharedentity.SubscriptionTier `json:"subscription_tier"`
}