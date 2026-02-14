package auth

import (
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

type Config struct {
	SignKey              string `koanf:"sign_key"`
	AccessExpirationTime time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject string `koanf:"access_subject"`
	RefreshSubject string `koanf:"refresh_subject"`
}

const (
	OpCreateAccessToken  = richerror.Op("authservice.CreateAccessToken")
	OpCreateRefreshToken = richerror.Op("authservice.CreateRefreshToken")
	OpParseToken         = richerror.Op("authservice.ParseToken")
	OpCreateToken        = richerror.Op("authservice.createToken")
)

type Service struct {
	config Config
}

func New(cfg Config) Service {
	return Service{
		config: cfg,
	}
}