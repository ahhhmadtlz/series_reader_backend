package config

import "github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"

func Default() Config {
	cfx := Config{
		// Auth: auth.Config {
		// 	AccessExpirationTime: AccessTokenExpireDuration,
		// 	RefreshExpirationTime: RefreshTokenExpireDuration,
		// 	AccessSubject: AccessTokenSubject,
		// 	RefreshSubject: RefreshTokenSubject,
		// },
		Logger: logger.Config{
			UseLocalTime:     LoggerUseLocalTime,
			FileMaxSizeInMB:  LoggerFileMaxSizeInMB,
			FileMaxAgeInDays: LoggerFileMaxAgeInDays,
			MaxBackups:       LoggerMaxBackups,
			Compress:         LoggerCompress,
		},
		// Redis: Redis{
		// 	Host:     "localhost",
		// 	Port:     6379,
		// 	Password: "",
		// 	DB:       0,
		// 	PoolSize: 10,
		// },
		// Jobs: Jobs{
		// 	Concurrency: 10,
		// 	Queues: map[string]int{
		// 		"default":  6,
		// 		"critical": 10,
		// 		"low":      1,
		// 	},
		// },
	}
	return cfx
}