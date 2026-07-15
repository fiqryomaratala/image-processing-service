package service

import (
	"bytes"
	"context"
	"time"

	filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/entity"
	imagerepository "github.com/fiqryomaratala/image-processing-service/backend/internal/image/repository"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	queuepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/queue"
	storagepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageService struct {
	repository imagerepository.Repository
	validator  uploadValidator
	storage    storagepkg.Storage
	publisher  queuepkg.Publisher
	logger     *zap.Logger
}

type uploadValidator interface {
	Validate(upload filepkg.Upload) (*filepkg.ValidationResult, error)
}

var _ Service = (*ImageService)(nil)

func NewImageService(repository imagerepository.Repository, validator uploadValidator, storage storagepkg.Storage, publisher queuepkg.Publisher) *ImageService {
	return &ImageService{
		repository: repository,
		validator:  validator,
		storage:    storage,
		publisher:  publisher,
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

	s.logger.Info("upload succeeded to minio",
		zap.String("object_key", uploadResult.ObjectKey),
		zap.String("bucket", uploadResult.Bucket),
		zap.String("content_type", uploadResult.ContentType),
		zap.Int64("size", uploadResult.Size),
	)

	image := buildImageEntity(result, uploadResult)

	if err := s.repository.Create(ctx, image); err != nil {
		s.logger.Warn("metadata save failed, starting compensation",
			zap.Error(err),
			zap.String("object_key", uploadResult.ObjectKey),
		)

		if deleteErr := s.storage.Delete(ctx, uploadResult.Bucket, uploadResult.ObjectKey); deleteErr != nil {
			s.logger.Warn("compensation failed",
				zap.Error(deleteErr),
				zap.String("object_key", uploadResult.ObjectKey),
			)
		} else {
			s.logger.Info("compensation succeeded",
				zap.String("object_key", uploadResult.ObjectKey),
			)
		}

		return nil, err
	}

	s.logger.Info("metadata saved",
		zap.String("image_id", image.ID.String()),
		zap.String("object_key", image.ObjectKey),
		zap.String("bucket", image.BucketName),
	)

	job := buildImageJob(image)

	s.logger.Info("job publish started",
		zap.String("image_id", job.ImageID),
		zap.String("object_key", job.ObjectKey),
	)

	if err := s.publisher.PublishImageJob(ctx, job); err != nil {
		s.logger.Warn("job publish failed",
			zap.Error(err),
			zap.String("image_id", job.ImageID),
			zap.String("object_key", job.ObjectKey),
		)

		return nil, err
	}

	if err := s.repository.UpdateStatus(ctx, image.ID, entity.StatusQueued); err != nil {
		s.logger.Warn("failed to update image status to queued",
			zap.Error(err),
			zap.String("image_id", image.ID.String()),
			zap.String("object_key", image.ObjectKey),
		)
		return nil, err
	}

	s.logger.Info("job publish succeeded",
		zap.String("image_id", image.ID.String()),
		zap.String("object_key", image.ObjectKey),
		zap.String("status", string(entity.StatusQueued)),
	)

	return &dto.UploadResponse{
		ID:        image.ID.String(),
		ObjectKey: image.ObjectKey,
		Status:    string(entity.StatusQueued),
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

func buildImageEntity(validationResult *filepkg.ValidationResult, uploadResult *storagepkg.UploadResult) *entity.Image {
	now := time.Now().UTC()

	return &entity.Image{
		ID:               uuid.New(),
		OriginalFilename: validationResult.OriginalFilename,
		StoredFilename:   validationResult.SanitizedFilename,
		ObjectKey:        uploadResult.ObjectKey,
		BucketName:       uploadResult.Bucket,
		ContentType:      uploadResult.ContentType,
		FileSize:         uploadResult.Size,
		Width:            validationResult.Image.Width,
		Height:           validationResult.Image.Height,
		Status:           entity.StatusUploaded,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func buildImageJob(image *entity.Image) queuepkg.ImageJob {
	return queuepkg.ImageJob{
		ImageID:     image.ID.String(),
		ObjectKey:   image.ObjectKey,
		BucketName:  image.BucketName,
		ContentType: image.ContentType,
		Status:      string(entity.StatusQueued),
		UploadedAt:  image.CreatedAt,
	}
}
