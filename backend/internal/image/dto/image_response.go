package dto

import (
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/entity"
)

type ImageResponse struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	StoredFilename   string    `json:"stored_filename"`
	ObjectKey        string    `json:"object_key"`
	BucketName       string    `json:"bucket_name"`
	ContentType      string    `json:"content_type"`
	FileSize         int64     `json:"file_size"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewImageResponse(image *entity.Image) *ImageResponse {
	if image == nil {
		return nil
	}

	return &ImageResponse{
		ID:               image.ID.String(),
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
