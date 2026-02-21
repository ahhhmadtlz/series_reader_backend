package config

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
)

type HTTPServer struct {
	Port int `koanf:"port"`
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
	BasePath string `koanf:"base_path"` 
	BaseURL  string `koanf:"base_url"`  
	MaxAvatarSizeMB int `koanf:"max_avatar_size_mb"` // 5
	MaxCoverSizeMB  int `koanf:"max_cover_size_mb"`  // 10
	MaxPageSizeMB   int `koanf:"max_page_size_mb"`   // 15
	AllowedMimeTypes []string `koanf:"allowed_mime_types"` // ["image/jpeg", "image/png", "image/webp"]
}

type Config struct {
	HTTPServer HTTPServer    `koanf:"http_server"`
  Postgres   Postgres      `koanf:"postgres"`
	Logger     logger.Config `koanf:"logger"`
  Auth       auth.Config   `koanf:"auth"`
  Upload     Upload        `koanf:"upload"`
}





