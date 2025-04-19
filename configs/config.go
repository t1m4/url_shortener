package configs

import "os"

type Config struct {
	ENVIRONMENT  string
	DOMAIN       string
	SERVER_HOST  string
	POSTGRES_DSN string
}

func LoadConfig() *Config {
	return &Config{
		ENVIRONMENT:  os.Getenv("ENVIRONMENT"),
		DOMAIN:       os.Getenv("DOMAIN"),
		SERVER_HOST:  os.Getenv("SERVER_HOST"),
		POSTGRES_DSN: os.Getenv("POSTGRES_DSN"),
	}
}
