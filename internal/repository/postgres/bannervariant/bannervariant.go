package bannervariant

import "database/sql"

type BannerVariantRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *BannerVariantRepo {
	return &BannerVariantRepo{db: db}
}