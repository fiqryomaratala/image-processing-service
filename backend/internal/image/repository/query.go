package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ImageRepository) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(&imageModel{})
}

func (r *ImageRepository) findByIDQuery(ctx context.Context, id uuid.UUID) *gorm.DB {
	return r.baseQuery(ctx).Where("id = ?", id)
}

func (r *ImageRepository) findByObjectKeyQuery(ctx context.Context, objectKey string) *gorm.DB {
	return r.baseQuery(ctx).Where("object_key = ?", objectKey)
}
