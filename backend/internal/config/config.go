package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

const (
	defaultAppEnv   = "development"
	defaultAppHost  = "0.0.0.0"
	defaultAppPort  = "8080"
	defaultLogLevel = "info"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	RabbitMQ RabbitMQConfig
	MinIO    MinIOConfig
}

type AppConfig struct {
	Env      string
	Host     string
	Port     string
	LogLevel string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type MinIOConfig struct {
	Host         string
	APIPort      string
	ConsolePort  string
	RootUser     string
	RootPassword string
}

var (
	cfg  *Config
	once sync.Once
	err  error
)

func Load() (*Config, error) {
	once.Do(func() {
		_ = godotenv.Load()

		c := &Config{
			App: AppConfig{
				Env:      getEnv("APP_ENV", defaultAppEnv),
				Host:     getEnv("APP_HOST", defaultAppHost),
				Port:     getEnv("APP_PORT", defaultAppPort),
				LogLevel: getEnv("LOG_LEVEL", defaultLogLevel),
			},
			Postgres: PostgresConfig{
				Host:     getEnv("POSTGRES_HOST", ""),
				Port:     getEnv("POSTGRES_PORT", ""),
				Name:     getEnv("POSTGRES_DB", ""),
				User:     getEnv("POSTGRES_USER", ""),
				Password: getEnv("POSTGRES_PASSWORD", ""),
			},
			RabbitMQ: RabbitMQConfig{
				Host:     getEnv("RABBITMQ_HOST", ""),
				Port:     getEnv("RABBITMQ_PORT", ""),
				User:     getEnv("RABBITMQ_USER", ""),
				Password: getEnv("RABBITMQ_PASSWORD", ""),
			},
			MinIO: MinIOConfig{
				Host:         getEnv("MINIO_HOST", ""),
				APIPort:      getEnv("MINIO_API_PORT", ""),
				ConsolePort:  getEnv("MINIO_CONSOLE_PORT", ""),
				RootUser:     getEnv("MINIO_ROOT_USER", ""),
				RootPassword: getEnv("MINIO_ROOT_PASSWORD", ""),
			},
		}

		err = c.Validate()
		if err != nil {
			return
		}

		cfg = c
	})

	return cfg, err
}

func MustLoad() *Config {
	c, err := Load()
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) Validate() error {
	required := map[string]string{
		"POSTGRES_HOST":         c.Postgres.Host,
		"POSTGRES_PORT":         c.Postgres.Port,
		"POSTGRES_DB":           c.Postgres.Name,
		"POSTGRES_USER":         c.Postgres.User,
		"POSTGRES_PASSWORD":     c.Postgres.Password,
		"RABBITMQ_HOST":         c.RabbitMQ.Host,
		"RABBITMQ_PORT":         c.RabbitMQ.Port,
		"RABBITMQ_USER":         c.RabbitMQ.User,
		"RABBITMQ_PASSWORD":     c.RabbitMQ.Password,
		"MINIO_HOST":            c.MinIO.Host,
		"MINIO_API_PORT":        c.MinIO.APIPort,
		"MINIO_CONSOLE_PORT":    c.MinIO.ConsolePort,
		"MINIO_ROOT_USER":       c.MinIO.RootUser,
		"MINIO_ROOT_PASSWORD":   c.MinIO.RootPassword,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("%s is required", key)
		}
	}

	if err := validatePort("APP_PORT", c.App.Port); err != nil {
		return err
	}
	if err := validatePort("POSTGRES_PORT", c.Postgres.Port); err != nil {
		return err
	}
	if err := validatePort("RABBITMQ_PORT", c.RabbitMQ.Port); err != nil {
		return err
	}
	if err := validatePort("MINIO_API_PORT", c.MinIO.APIPort); err != nil {
		return err
	}
	if err := validatePort("MINIO_CONSOLE_PORT", c.MinIO.ConsolePort); err != nil {
		return err
	}

	return nil
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

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s must be a valid number", name)
	}

	if port < 1 || port > 65535 {
		return errors.New(name + " must be between 1 and 65535")
	}

	return nil
}