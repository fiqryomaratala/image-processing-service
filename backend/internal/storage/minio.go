package storage

import (
	"fmt"
	"sync"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	client *minio.Client
	mu     sync.Mutex
)

func Initialize() error {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		return nil
	}

	cfg := config.Get()
	log := logger.Get()

	log.Info("Connecting to MinIO...", zap.String("endpoint", buildEndpoint(cfg.MinIO)), zap.Bool("use_ssl", cfg.MinIO.UseSSL))

	minioClient, err := minio.New(buildEndpoint(cfg.MinIO), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.RootUser, cfg.MinIO.RootPassword, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		log.Error("Failed to connect MinIO", zap.Error(err))
		return fmt.Errorf("failed to connect minio: %w", err)
	}

	client = minioClient
	log.Info("MinIO connected successfully")

	return nil
}

func GetClient() *minio.Client {
	if client == nil {
		panic("minio is not initialized: call storage.Initialize() before storage.GetClient()")
	}

	return client
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	client = nil

	return nil
}

func buildEndpoint(cfg config.MinIOConfig) string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.APIPort)
}
