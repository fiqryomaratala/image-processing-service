package queue

import (
	"time"
)

const ImageProcessingQueue = "image.processing"

type ImageJob struct {
	ImageID     string    `json:"image_id"`
	ObjectKey   string    `json:"object_key"`
	BucketName  string    `json:"bucket_name"`
	ContentType string    `json:"content_type"`
	Status      string    `json:"status"`
	UploadedAt  time.Time `json:"uploaded_at"`
}
