package storage

import (
	"context"
	"fmt"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
)

func Health() error {
	minioClient := GetClient()
	if minioClient == nil {
		return fmt.Errorf("minio client is not initialized")
	}

	exists, err := minioClient.BucketExists(context.Background(), config.Get().MinIO.BucketName)
	if err != nil {
		return fmt.Errorf("minio health check failed: %w", err)
	}

	if !exists {
		return fmt.Errorf("minio health check failed: bucket %s is not accessible", config.Get().MinIO.BucketName)
	}

	return nil
}
