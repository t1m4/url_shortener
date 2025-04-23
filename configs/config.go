package configs

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type db struct {
	PostgresDsn string `yaml:"PostgresDsn"`
}
type App struct {
	Domain     string `yaml:"Domain"`
	ServerHost string `yaml:"ServerHost"`
}
type logger struct {
	Level string `yaml:"Level"`
}

type Limiter struct {
	Duration   time.Duration `yaml:"Duration"`
	EventCount int           `yaml:"EventCount"`
	Burst      int           `yaml:"Burst"`
}
type RateLimiterConfig struct {
	CleaningPeriod time.Duration `yaml:"CleaningPeriod"`
	ExpiresPeriod  time.Duration `yaml:"ExpiresPeriod"`
	Limiters       []Limiter     `yaml:"Limiters"`
}
type Config struct {
	Environment string            `yaml:"Environment"`
	App         App               `yaml:"App"`
	Db          db                `yaml:"Db"`
	Logger      logger            `yaml:"Logger"`
	RateLimiter RateLimiterConfig `yaml:"RateLimiter"`
}

func LoadConfig() (*Config, error) {
	var configFileNameByEnv = map[string]string{
		Dev:     "configs/local.yaml",
		Staging: "config.yaml",
	}
	env := os.Getenv("ENVIRONMENT")
	data, err := os.ReadFile(configFileNameByEnv[env])
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
