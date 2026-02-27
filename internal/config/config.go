package config

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
)

type HTTPServer struct {
	Port int `koanf:"port"`
	BodyLimitMB int `koanf:"body_limit_mb"`
}

type CORS struct {
	AllowedOrigins []string `koanf:"allowed_origins"`
}


type Postgres struct {
    Host     string `koanf:"host"`
    Port     int    `koanf:"port"`
    Username string `koanf:"username"`
    Password string `koanf:"password"`
    DBName   string `koanf:"db_name"`
    SSLMode  string `koanf:"ssl_mode"`
}

type Upload struct {
	BasePath         string   `koanf:"base_path"`
	BaseURL          string   `koanf:"base_url"`
	MaxAvatarSizeMB  int      `koanf:"max_avatar_size_mb"`
	MaxCoverSizeMB   int      `koanf:"max_cover_size_mb"`
	MaxPageSizeMB    int      `koanf:"max_page_size_mb"`
	MaxBannerSizeMB  int      `koanf:"max_banner_size_mb"`
	MaxThumbnailSizeMB int    `koanf:"max_thumbnail_size_mb"`
	AllowedMimeTypes []string `koanf:"allowed_mime_types"`
}

type Config struct {
	HTTPServer HTTPServer    `koanf:"http_server"`
  Postgres   Postgres      `koanf:"postgres"`
	Logger     logger.Config `koanf:"logger"`
  Auth       auth.Config   `koanf:"auth"`
  Upload     Upload        `koanf:"upload"`
	CORS       CORS          `koanf:"cors"`
}





