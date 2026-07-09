package service

import (
	"context"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, request dto.UploadRequest) (*dto.UploadResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ImageResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
