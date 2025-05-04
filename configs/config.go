package configs

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type db struct {
	PostgresDsn     string        `yaml:"PostgresDsn" json:"PostgresDsn"`
	MaxIdleConns    int           `yaml:"MaxIdleConns" json:"MaxIdleConns"`
	MaxOpenConns    int           `yaml:"MaxOpenConns" json:"MaxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"ConnMaxLifetime" json:"ConnMaxLifetime"`
}
type App struct {
	Domain       string        `yaml:"Domain" json:"Domain"`
	ServerHost   string        `yaml:"ServerHost" json:"ServerHost"`
	ReadTimeout  time.Duration `yaml:"ReadTimeout" json:"ReadTimeout"`
	WriteTimeout time.Duration `yaml:"WriteTimeout" json:"WriteTimeout"`
	IdleTimeout  time.Duration `yaml:"IdleTimeout" json:"IdleTimeout"`
}
type logger struct {
	Level string `yaml:"Level" json:"Level"`
}

type Limiter struct {
	Duration   time.Duration `yaml:"Duration" json:"Duration"`
	EventCount int           `yaml:"EventCount" json:"EventCount"`
	Burst      int           `yaml:"Burst" json:"Burst"`
}
type RateLimiterConfig struct {
	CleaningPeriod time.Duration `yaml:"CleaningPeriod" json:"CleaningPeriod"`
	ExpiresPeriod  time.Duration `yaml:"ExpiresPeriod" json:"ExpiresPeriod"`
	Limiters       []Limiter     `yaml:"Limiters" json:"Limiters"`
}
type Config struct {
	Environment string            `yaml:"Environment" json:"Environment"`
	App         App               `yaml:"App" json:"App"`
	Db          db                `yaml:"Db" json:"Db"`
	Logger      logger            `yaml:"Logger" json:"Logger"`
	RateLimiter RateLimiterConfig `yaml:"RateLimiter" json:"RateLimiter"`
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
