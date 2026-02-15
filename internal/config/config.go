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

type Config struct {
	HTTPServer HTTPServer    `koanf:"http_server"`
  Postgres   Postgres      `koanf:"postgres"`
	Logger     logger.Config `koanf:"logger"`
  Auth       auth.Config   `koanf:"auth"`
}