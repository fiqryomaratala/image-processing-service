package logger

import (
	"strings"
	"sync"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
	initErr  error
)

func New(cfg config.LoggerConfig, app config.AppConfig) (*zap.Logger, error) {
	once.Do(func() {
		zapConfig := buildConfig(cfg, app)

		instance, initErr = zapConfig.Build()
	})

	return instance, initErr
}

func MustNew(cfg config.LoggerConfig, app config.AppConfig) *zap.Logger {
	logger, err := New(cfg, app)
	if err != nil {
		panic(err)
	}

	return logger
}

func buildConfig(cfg config.LoggerConfig, app config.AppConfig) zap.Config {
	zapConfig := zap.NewProductionConfig()

	if strings.EqualFold(app.Env, "development") {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.Level = zap.NewAtomicLevelAt(parseLevel(cfg.Level))
	zapConfig.EncoderConfig.TimeKey = "timestamp"

	return zapConfig
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
