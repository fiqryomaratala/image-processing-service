package service

import (
	"context"

	filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
	imagerepository "github.com/fiqryomaratala/image-processing-service/backend/internal/image/repository"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageService struct {
	repository imagerepository.Repository
	validator  uploadValidator
	logger     *zap.Logger
}

type uploadValidator interface {
	Validate(upload filepkg.Upload) (*filepkg.ValidationResult, error)
}

var _ Service = (*ImageService)(nil)

func NewImageService(repository imagerepository.Repository, validator uploadValidator) *ImageService {
	return &ImageService{
		repository: repository,
		validator:  validator,
		logger:     resolveLogger(),
	}
}

func (s *ImageService) Upload(ctx context.Context, request dto.UploadRequest) (*dto.UploadResponse, error) {
	if err := validateUploadRequest(request); err != nil {
		s.logger.Warn("validation failed", zap.Error(err))
		return nil, err
	}

	result, err := s.validator.Validate(request.File)
	if err != nil {
		s.logger.Warn("validation failed", zap.Error(err))
		return nil, err
	}

	s.logger.Info("validation passed",
		zap.String("filename", result.SanitizedFilename),
		zap.String("content_type", result.MIMEType),
		zap.Int64("size", result.Size),
	)

	return &dto.UploadResponse{
		Filename:    result.SanitizedFilename,
		ContentType: result.MIMEType,
		Size:        result.Size,
	}, nil
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

func resolveLogger() *zap.Logger {
	log := zap.NewNop()

	defer func() {
		_ = recover()
	}()

	log = logger.Get()

	return log
}
