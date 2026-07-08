package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/minio/minio-go/v7"
)

func UploadObject(ctx context.Context, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	info, err := GetClient().PutObject(ctx, config.Get().MinIO.BucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("failed to upload object %s: %w", objectName, err)
	}

	return info, nil
}

func DownloadObject(ctx context.Context, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	object, err := GetClient().GetObject(ctx, config.Get().MinIO.BucketName, objectName, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to download object %s: %w", objectName, err)
	}

	return object, nil
}

func DeleteObject(ctx context.Context, objectName string, opts minio.RemoveObjectOptions) error {
	if err := GetClient().RemoveObject(ctx, config.Get().MinIO.BucketName, objectName, opts); err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectName, err)
	}

	return nil
}

func StatObject(ctx context.Context, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	info, err := GetClient().StatObject(ctx, config.Get().MinIO.BucketName, objectName, opts)
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to stat object %s: %w", objectName, err)
	}

	return info, nil
}
