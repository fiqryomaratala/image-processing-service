package repository

import (
	"context"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/entity"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, image *entity.Image) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	FindByObjectKey(ctx context.Context, objectKey string) (*entity.Image, error)
	ExistsByObjectKey(ctx context.Context, objectKey string) (bool, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.Status) error
	Delete(ctx context.Context, id uuid.UUID) error
}
