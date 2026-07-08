package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/entity"
	imageerrors "github.com/fiqryomaratala/image-processing-service/backend/internal/image/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

type imageModel struct {
	ID               uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	OriginalFilename string    `gorm:"column:original_filename"`
	StoredFilename   string    `gorm:"column:stored_filename"`
	ObjectKey        string    `gorm:"column:object_key"`
	BucketName       string    `gorm:"column:bucket_name"`
	ContentType      string    `gorm:"column:content_type"`
	FileSize         int64     `gorm:"column:file_size"`
	Width            int       `gorm:"column:width"`
	Height           int       `gorm:"column:height"`
	Status           string    `gorm:"column:status"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (imageModel) TableName() string {
	return "images"
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

func (r *ImageRepository) Create(ctx context.Context, image *entity.Image) error {
	model := newImageModel(image)

	if err := r.baseQuery(ctx).Create(model).Error; err != nil {
		return wrapDatabaseError("create image", err)
	}

	assignEntityFromModel(image, model)

	return nil
}

func (r *ImageRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	var model imageModel
	if err := r.findByIDQuery(ctx, id).First(&model).Error; err != nil {
		return nil, mapQueryError("find image by id", err)
	}

	return model.toEntity(), nil
}

func (r *ImageRepository) FindByObjectKey(ctx context.Context, objectKey string) (*entity.Image, error) {
	var model imageModel
	if err := r.findByObjectKeyQuery(ctx, objectKey).First(&model).Error; err != nil {
		return nil, mapQueryError("find image by object key", err)
	}

	return model.toEntity(), nil
}

func (r *ImageRepository) ExistsByObjectKey(ctx context.Context, objectKey string) (bool, error) {
	var count int64
	if err := r.findByObjectKeyQuery(ctx, objectKey).Count(&count).Error; err != nil {
		return false, wrapDatabaseError("check image object key existence", err)
	}

	return count > 0, nil
}

func (r *ImageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.Status) error {
	result := r.findByIDQuery(ctx, id).Updates(map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	})
	if result.Error != nil {
		return wrapDatabaseError("update image status", result.Error)
	}

	if result.RowsAffected == 0 {
		return imageerrors.ImageNotFound()
	}

	return nil
}

func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.findByIDQuery(ctx, id).Delete(&imageModel{})
	if result.Error != nil {
		return wrapDatabaseError("delete image", result.Error)
	}

	if result.RowsAffected == 0 {
		return imageerrors.ImageNotFound()
	}

	return nil
}

func newImageModel(image *entity.Image) *imageModel {
	if image == nil {
		return &imageModel{}
	}

	return &imageModel{
		ID:               image.ID,
		OriginalFilename: image.OriginalFilename,
		StoredFilename:   image.StoredFilename,
		ObjectKey:        image.ObjectKey,
		BucketName:       image.BucketName,
		ContentType:      image.ContentType,
		FileSize:         image.FileSize,
		Width:            image.Width,
		Height:           image.Height,
		Status:           string(image.Status),
		CreatedAt:        image.CreatedAt,
		UpdatedAt:        image.UpdatedAt,
	}
}

func (m *imageModel) toEntity() *entity.Image {
	if m == nil {
		return nil
	}

	return &entity.Image{
		ID:               m.ID,
		OriginalFilename: m.OriginalFilename,
		StoredFilename:   m.StoredFilename,
		ObjectKey:        m.ObjectKey,
		BucketName:       m.BucketName,
		ContentType:      m.ContentType,
		FileSize:         m.FileSize,
		Width:            m.Width,
		Height:           m.Height,
		Status:           entity.Status(m.Status),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func assignEntityFromModel(image *entity.Image, model *imageModel) {
	if image == nil || model == nil {
		return
	}

	*image = *model.toEntity()
}

func mapQueryError(action string, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return imageerrors.ImageNotFound().WithCause(err)
	}

	return wrapDatabaseError(action, err)
}

func wrapDatabaseError(action string, err error) error {
	return apperrors.Internal(fmt.Sprintf("failed to %s", action)).WithCause(err)
}
