package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	instance *gorm.DB
	mu       sync.Mutex
)

func Initialize() error {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		return nil
	}

	cfg := config.Get()
	log := logger.Get()

	log.Info("Connecting to PostgreSQL...", zap.String("host", cfg.Database.Host), zap.String("port", cfg.Database.Port), zap.String("database", cfg.Database.Database))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  buildDSN(cfg.Database),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 newGORMLogger(log, cfg.App.Env),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		log.Error("Failed to connect PostgreSQL", zap.Error(err))
		return fmt.Errorf("failed to connect postgres: %w", err)
	}

	instance = db
	log.Info("PostgreSQL connected successfully")

	return nil
}

func Get() *gorm.DB {
	if instance == nil {
		panic("database is not initialized: call database.Initialize() before database.Get()")
	}

	return instance
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		return nil
	}

	sqlDB, err := instance.DB()
	if err != nil {
		return fmt.Errorf("failed to access postgres sql db: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close postgres connection: %w", err)
	}

	instance = nil

	return nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
		cfg.TimeZone,
	)
}

func newGORMLogger(log *zap.Logger, appEnv string) gormlogger.Interface {
	logLevel := gormlogger.Warn
	if appEnv == "development" {
		logLevel = gormlogger.Info
	}

	return gormlogger.New(
		&gormLogWriter{logger: log.Named("gorm")},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  appEnv == "development",
		},
	)
}

type gormLogWriter struct {
	logger *zap.Logger
}

func (w *gormLogWriter) Printf(format string, args ...any) {
	w.logger.Debug(fmt.Sprintf(format, args...))
}
