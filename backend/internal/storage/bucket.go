package storage

import (
	"context"
	"fmt"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

const defaultBucketRegion = "us-east-1"

func EnsureBucket() error {
	cfg := config.Get()
	log := logger.Get()
	minioClient := GetClient()

	exists, err := minioClient.BucketExists(context.Background(), cfg.MinIO.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check minio bucket %s: %w", cfg.MinIO.BucketName, err)
	}

	if exists {
		log.Info("Bucket is ready", zap.String("bucket", cfg.MinIO.BucketName))
		log.Info("Bucket already exists", zap.String("bucket", cfg.MinIO.BucketName))
		return nil
	}

	if err := minioClient.MakeBucket(context.Background(), cfg.MinIO.BucketName, minio.MakeBucketOptions{
		Region: defaultBucketRegion,
	}); err != nil {
		return fmt.Errorf("failed to create minio bucket %s: %w", cfg.MinIO.BucketName, err)
	}

	log.Info("Bucket created successfully", zap.String("bucket", cfg.MinIO.BucketName))
	log.Info("Bucket is ready", zap.String("bucket", cfg.MinIO.BucketName))

	return nil
}
