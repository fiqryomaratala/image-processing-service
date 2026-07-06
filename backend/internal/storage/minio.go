package storage

import (
	"fmt"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIO(cfg config.MinIOConfig) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.APIPort)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.RootPassword, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	return client, nil
}