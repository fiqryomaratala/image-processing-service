package dto

import (
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/entity"
)

type JobResponse struct {
	ImageID   string    `json:"image_id"`
	ObjectKey string    `json:"object_key"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewJobResponse(image *entity.Image) *JobResponse {
	if image == nil {
		return nil
	}

	return &JobResponse{
		ImageID:   image.ID.String(),
		ObjectKey: image.ObjectKey,
		Status:    string(image.Status),
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}
}
