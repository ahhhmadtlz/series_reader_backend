package config

func Default() Config {
	cfx := Config{
		// Auth: auth.Config {
		// 	AccessExpirationTime: AccessTokenExpireDuration,
		// 	RefreshExpirationTime: RefreshTokenExpireDuration,
		// 	AccessSubject: AccessTokenSubject,
		// 	RefreshSubject: RefreshTokenSubject,
		// },
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
		Postgres: Postgres{
    Host:     "localhost",
    Port:     5432,
    Username: "postgres",
    Password: "postgres",
    DBName:   "series_reader",
    SSLMode:  "disable",
},
	}
	return cfx
}