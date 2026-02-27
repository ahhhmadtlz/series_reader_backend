package config

func Default() Config {
	return Config{
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "", // set via SERIES_POSTGRES__PASSWORD env var
			DBName:   "series_reader",
			SSLMode:  "disable",
		},
		Upload: Upload{
			BasePath:           "./uploads",
			BaseURL:            "http://localhost:8080/uploads",
			MaxAvatarSizeMB:    5,
			MaxCoverSizeMB:     10,
			MaxPageSizeMB:      15,
			MaxBannerSizeMB:    10,
			MaxThumbnailSizeMB: 5,
			AllowedMimeTypes:   []string{"image/jpeg", "image/jpg", "image/png", "image/webp"},
		},
	}
}