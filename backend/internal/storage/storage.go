package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, object Object) (*UploadResult, error)
	Delete(ctx context.Context, bucket, objectKey string) error
	Exists(ctx context.Context, bucket, objectKey string) (bool, error)
	GetObjectURL(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
}

type Object struct {
	Bucket      string
	ObjectKey   string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type UploadResult struct {
	Bucket      string
	ObjectKey   string
	ContentType string
	Size        int64
	ETag        string
}
