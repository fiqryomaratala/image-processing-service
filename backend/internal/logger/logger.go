package logger

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	initOnce sync.Once
	initErr  error
)

func Initialize() error {
	initOnce.Do(func() {
		cfg := config.Get()
		zapConfig, err := buildConfig(cfg.Logger, cfg.App)
		if err != nil {
			initErr = err
			return
		}

		instance, initErr = zapConfig.Build()
	})

	return initErr
}

func Get() *zap.Logger {
	if instance == nil {
		panic("logger is not initialized: call logger.Initialize() before logger.Get()")
	}

	return instance
}

func buildConfig(cfg config.LoggerConfig, app config.AppConfig) (zap.Config, error) {
	zapConfig := zap.NewProductionConfig()

	if strings.EqualFold(app.Env, "development") {
		zapConfig = zap.NewDevelopmentConfig()
	}

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return zap.Config{}, err
	}

	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.DisableStacktrace = true

	return zapConfig, nil
}

func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info", "":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid logger level: %s", level)
	}
}
