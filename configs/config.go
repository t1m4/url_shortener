package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Db struct {
	POSTGRES_DSN string `yaml:"postgres_dsn"`
}
type App struct {
	DOMAIN      string `yaml:"domain"`
	SERVER_HOST string `yaml:"server_host"`
}
type Logger struct {
	LEVEL string `yaml:"level"`
}
type Config struct {
	ENVIRONMENT string `yaml:"environment"`
	APP         App    `yaml:"app"`
	DB          Db     `yaml:"db"`
	LOGGER      Logger `yaml:"logger"`
}

func LoadConfig() (*Config, error) {
	var configFileNameByEnv = map[string]string{
		DEV:     "configs/local.yaml",
		STAGING: "configs/stage.yaml",
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
