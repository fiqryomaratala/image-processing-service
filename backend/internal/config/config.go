package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	defaultAppEnv      = "development"
	defaultAppHost     = "0.0.0.0"
	defaultAppPort     = "8080"
	defaultLogLevel    = "info"
	defaultDBSSLMode   = "disable"
	defaultDBTimeZone  = "UTC"
	defaultMinIOUseSSL = false
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	MinIO    MinIOConfig
	Logger   LoggerConfig
}

type AppConfig struct {
	Name string
	Env  string
	Host string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
	TimeZone string
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
	BucketName   string
	UseSSL       bool
}

type LoggerConfig struct {
	Level string
}

var (
	instance *Config
	loadOnce sync.Once
	loadErr  error
)

func Load() error {
	loadOnce.Do(func() {
		v := newViper()

		if err := loadEnvFiles(v); err != nil {
			loadErr = err
			return
		}

		cfg := &Config{
			App: AppConfig{
				Name: v.GetString("APP_NAME"),
				Env:  v.GetString("APP_ENV"),
				Host: v.GetString("APP_HOST"),
				Port: v.GetString("APP_PORT"),
			},
			Database: DatabaseConfig{
				Host:     v.GetString("POSTGRES_HOST"),
				Port:     v.GetString("POSTGRES_PORT"),
				Database: v.GetString("POSTGRES_DB"),
				User:     v.GetString("POSTGRES_USER"),
				Password: v.GetString("POSTGRES_PASSWORD"),
				SSLMode:  v.GetString("POSTGRES_SSLMODE"),
				TimeZone: v.GetString("POSTGRES_TIMEZONE"),
			},
			RabbitMQ: RabbitMQConfig{
				Host:     v.GetString("RABBITMQ_HOST"),
				Port:     v.GetString("RABBITMQ_PORT"),
				User:     v.GetString("RABBITMQ_USER"),
				Password: v.GetString("RABBITMQ_PASSWORD"),
			},
			MinIO: MinIOConfig{
				Host:         v.GetString("MINIO_HOST"),
				APIPort:      v.GetString("MINIO_API_PORT"),
				ConsolePort:  v.GetString("MINIO_CONSOLE_PORT"),
				RootUser:     v.GetString("MINIO_ROOT_USER"),
				RootPassword: v.GetString("MINIO_ROOT_PASSWORD"),
				BucketName:   v.GetString("MINIO_BUCKET_NAME"),
				UseSSL:       v.GetBool("MINIO_USE_SSL"),
			},
			Logger: LoggerConfig{
				Level: v.GetString("LOG_LEVEL"),
			},
		}

		if err := cfg.validate(); err != nil {
			loadErr = err
			return
		}

		instance = cfg
	})

	return loadErr
}

func Get() *Config {
	if instance == nil {
		panic("config is not loaded: call config.Load() before config.Get()")
	}

	return instance
}

func (c AppConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("APP_NAME", "image-processing-service")
	v.SetDefault("APP_ENV", defaultAppEnv)
	v.SetDefault("APP_HOST", defaultAppHost)
	v.SetDefault("APP_PORT", defaultAppPort)
	v.SetDefault("LOG_LEVEL", defaultLogLevel)
	v.SetDefault("POSTGRES_SSLMODE", defaultDBSSLMode)
	v.SetDefault("POSTGRES_TIMEZONE", defaultDBTimeZone)
	v.SetDefault("MINIO_USE_SSL", defaultMinIOUseSSL)

	return v
}

func loadEnvFiles(v *viper.Viper) error {
	for _, filename := range []string{".env", ".env.local"} {
		if err := mergeEnvFile(v, filename); err != nil {
			return err
		}
	}

	return nil
}

func mergeEnvFile(v *viper.Viper, filename string) error {
	for _, basePath := range []string{".", ".."} {
		path := filepath.Join(basePath, filename)

		v.SetConfigFile(path)
		if err := v.MergeInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				continue
			}

			return fmt.Errorf("failed to load %s: %w", path, err)
		}

		return nil
	}

	return nil
}

func (c *Config) validate() error {
	required := map[string]string{
		"POSTGRES_HOST":       c.Database.Host,
		"POSTGRES_PORT":       c.Database.Port,
		"POSTGRES_DB":         c.Database.Database,
		"POSTGRES_USER":       c.Database.User,
		"POSTGRES_PASSWORD":   c.Database.Password,
		"POSTGRES_SSLMODE":    c.Database.SSLMode,
		"POSTGRES_TIMEZONE":   c.Database.TimeZone,
		"RABBITMQ_HOST":       c.RabbitMQ.Host,
		"RABBITMQ_PORT":       c.RabbitMQ.Port,
		"RABBITMQ_USER":       c.RabbitMQ.User,
		"RABBITMQ_PASSWORD":   c.RabbitMQ.Password,
		"MINIO_HOST":          c.MinIO.Host,
		"MINIO_API_PORT":      c.MinIO.APIPort,
		"MINIO_CONSOLE_PORT":  c.MinIO.ConsolePort,
		"MINIO_ROOT_USER":     c.MinIO.RootUser,
		"MINIO_ROOT_PASSWORD": c.MinIO.RootPassword,
		"MINIO_BUCKET_NAME":   c.MinIO.BucketName,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("configuration error: %s is required", key)
		}
	}

	for _, port := range []struct {
		name  string
		value string
	}{
		{name: "APP_PORT", value: c.App.Port},
		{name: "POSTGRES_PORT", value: c.Database.Port},
		{name: "RABBITMQ_PORT", value: c.RabbitMQ.Port},
		{name: "MINIO_API_PORT", value: c.MinIO.APIPort},
		{name: "MINIO_CONSOLE_PORT", value: c.MinIO.ConsolePort},
	} {
		if err := validatePort(port.name, port.value); err != nil {
			return err
		}
	}

	return nil
}

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("configuration error: %s must be a valid number", name)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("configuration error: %s must be between 1 and 65535", name)
	}

	return nil
}
