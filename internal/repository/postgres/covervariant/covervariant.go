package covervariant

import "database/sql"

type CoverVariantRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *CoverVariantRepo {
	return &CoverVariantRepo{db: db}
}