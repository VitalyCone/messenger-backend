package service


type Config struct {
	DatabaseURL     string
	RedisURL string
	TokenKey string
}

func NewConfig(databaseURL, redisURL, tokenKey string) *Config{
	return &Config{
		DatabaseURL : databaseURL,
		RedisURL: redisURL,
		TokenKey: tokenKey,
	}
}