package config

import (
	"fmt"
	"os"
)

const (
	defaultAppHost  = "0.0.0.0"
	defaultAppPort  = "8080"
	defaultAppEnv   = "development"
	defaultLogLevel = "info"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Env      string
	Host     string
	Port     string
	LogLevel string
}

func Load() Config {
	return Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", defaultAppEnv),
			Host:     getEnv("APP_HOST", defaultAppHost),
			Port:     getEnv("APP_PORT", defaultAppPort),
			LogLevel: getEnv("LOG_LEVEL", defaultLogLevel),
		},
	}
}

func (a AppConfig) Address() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
