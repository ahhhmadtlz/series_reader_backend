package config

import "github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer    `koanf:"http_server"`
	Logger     logger.Config `koanf:"logger"`
}