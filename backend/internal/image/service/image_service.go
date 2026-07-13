package service

import (
	"bytes"
	"context"

	filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
	imagerepository "github.com/fiqryomaratala/image-processing-service/backend/internal/image/repository"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	storagepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageService struct {
	repository imagerepository.Repository
	validator  uploadValidator
	storage    storagepkg.Storage
	logger     *zap.Logger
}

type uploadValidator interface {
	Validate(upload filepkg.Upload) (*filepkg.ValidationResult, error)
}

var _ Service = (*ImageService)(nil)

func NewImageService(repository imagerepository.Repository, validator uploadValidator, storage storagepkg.Storage) *ImageService {
	return &ImageService{
		repository: repository,
		validator:  validator,
		storage:    storage,
		logger:     resolveLogger(),
	}
}

func (s *ImageService) Upload(ctx context.Context, request dto.UploadRequest) (*dto.UploadResponse, error) {
	if err := validateUploadRequest(request); err != nil {
		s.logger.Warn("upload failed during validation", zap.Error(err))
		return nil, err
	}

	s.logger.Info("upload started",
		zap.String("filename", request.File.Filename),
		zap.Int64("size", request.File.Size),
	)

	result, err := s.validator.Validate(request.File)
	if err != nil {
		s.logger.Warn("upload failed during validation", zap.Error(err))
		return nil, err
	}

	uploadResult, err := s.storage.Upload(ctx, storagepkg.Object{
		ObjectKey:   result.ObjectKey,
		ContentType: result.MIMEType,
		Size:        result.Size,
		Reader:      bytes.NewReader(result.Content),
	})
	if err != nil {
		s.logger.Warn("upload failed", zap.Error(err), zap.String("object_key", result.ObjectKey))
		return nil, err
	}

	s.logger.Info("upload succeeded",
		zap.String("object_key", uploadResult.ObjectKey),
		zap.String("bucket", uploadResult.Bucket),
		zap.String("content_type", uploadResult.ContentType),
		zap.Int64("size", uploadResult.Size),
	)

	return &dto.UploadResponse{
		ObjectKey:   uploadResult.ObjectKey,
		Bucket:      uploadResult.Bucket,
		ContentType: uploadResult.ContentType,
		Size:        uploadResult.Size,
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
