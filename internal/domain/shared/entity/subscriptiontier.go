
package entity

type SubscriptionTier uint8


const (
	FreeTier SubscriptionTier = iota + 1
	PremiumTier
)

const (
	FreeTierStr = "free"
	PremiumTierStr = "premium"
)


func (s SubscriptionTier) String() string {
	switch s {
	case FreeTier:
		return FreeTierStr
	case PremiumTier:
		return PremiumTierStr
	}
	return  ""
}

func MapToSubscriptionTier(tierStr string) SubscriptionTier{
	switch tierStr {
		case FreeTierStr:
			return FreeTier
		case PremiumTierStr:
			return PremiumTier
	}
	return SubscriptionTier(0)
}


