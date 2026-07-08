package service

import (
	"context"
	"strings"
	"time"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/entity"
	imagerepository "github.com/fiqryomaratala/image-processing-service/backend/internal/image/repository"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageService struct {
	repository imagerepository.Repository
	logger     *zap.Logger
}

var _ Service = (*ImageService)(nil)

func NewImageService(repository imagerepository.Repository) *ImageService {
	return &ImageService{
		repository: repository,
		logger:     resolveLogger(),
	}
}

func (s *ImageService) Upload(ctx context.Context, request dto.UploadRequest) (*dto.JobResponse, error) {
	if err := validateUploadRequest(request); err != nil {
		return nil, err
	}

	exists, err := s.repository.ExistsByObjectKey(ctx, request.ObjectKey)
	if err != nil {
		return nil, err
	}

	if exists {
		s.logger.Warn("image upload rejected because object key already exists",
			zap.String("object_key", request.ObjectKey),
		)

		return nil, apperrors.Conflict("image object key already exists")
	}

	image := buildImageEntity(request)
	if err := s.repository.Create(ctx, image); err != nil {
		return nil, err
	}

	s.logger.Info("image metadata created",
		zap.String("image_id", image.ID.String()),
		zap.String("object_key", image.ObjectKey),
	)

	return dto.NewJobResponse(image), nil
}

func (s *ImageService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ImageResponse, error) {
	image, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.NewImageResponse(image), nil
}

func (s *ImageService) Delete(ctx context.Context, id uuid.UUID) error {
	image, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.Info("image metadata deleted",
		zap.String("image_id", id.String()),
		zap.String("object_key", image.ObjectKey),
	)

	return nil
}

func buildImageEntity(request dto.UploadRequest) *entity.Image {
	now := time.Now().UTC()

	return &entity.Image{
		ID:               uuid.New(),
		OriginalFilename: strings.TrimSpace(request.OriginalFilename),
		StoredFilename:   strings.TrimSpace(request.StoredFilename),
		ObjectKey:        strings.TrimSpace(request.ObjectKey),
		BucketName:       strings.TrimSpace(request.BucketName),
		ContentType:      strings.TrimSpace(request.ContentType),
		FileSize:         request.FileSize,
		Width:            request.Width,
		Height:           request.Height,
		Status:           entity.StatusUploaded,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func resolveLogger() *zap.Logger {
	log := zap.NewNop()

	defer func() {
		_ = recover()
	}()

	log = logger.Get()

	return log
}
