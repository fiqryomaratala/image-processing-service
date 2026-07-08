package entity

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusUploaded   Status = "uploaded"
	StatusQueued     Status = "queued"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusUploaded, StatusQueued, StatusProcessing, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

type Image struct {
	ID               uuid.UUID
	OriginalFilename string
	StoredFilename   string
	ObjectKey        string
	BucketName       string
	ContentType      string
	FileSize         int64
	Width            int
	Height           int
	Status           Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
