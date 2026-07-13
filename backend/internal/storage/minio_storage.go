package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

var _ Storage = (*MinIOStorage)(nil)

func NewMinIOStorage(client *minio.Client, bucket string) *MinIOStorage {
	return &MinIOStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *MinIOStorage) Upload(ctx context.Context, object Object) (*UploadResult, error) {
	client, err := s.ensureClient()
	if err != nil {
		return nil, err
	}

	bucket, err := s.resolveBucket(object.Bucket)
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, StorageUnavailable(fmt.Sprintf("failed to check bucket %s: %v", bucket, err))
	}
	if !exists {
		return nil, BucketUnavailable(bucket)
	}

	objectKey := normalizeObjectKey(object.ObjectKey)
	if objectKey == "" {
		return nil, UploadFailed("object key is required")
	}

	info, err := client.PutObject(ctx, bucket, objectKey, object.Reader, object.Size, minio.PutObjectOptions{
		ContentType: object.ContentType,
	})
	if err != nil {
		return nil, UploadFailed(fmt.Sprintf("failed to upload object %s: %v", objectKey, err))
	}

	return &UploadResult{
		Bucket:      bucket,
		ObjectKey:   objectKey,
		ContentType: object.ContentType,
		Size:        object.Size,
		ETag:        info.ETag,
	}, nil
}

func (s *MinIOStorage) Delete(ctx context.Context, bucket, objectKey string) error {
	client, err := s.ensureClient()
	if err != nil {
		return err
	}

	bucket, err = s.resolveBucket(bucket)
	if err != nil {
		return err
	}
	objectKey = normalizeObjectKey(objectKey)
	if objectKey == "" {
		return UploadFailed("object key is required")
	}

	if err := client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return UploadFailed(fmt.Sprintf("failed to delete object %s: %v", objectKey, err))
	}

	return nil
}

func (s *MinIOStorage) Exists(ctx context.Context, bucket, objectKey string) (bool, error) {
	client, err := s.ensureClient()
	if err != nil {
		return false, err
	}

	bucket, err = s.resolveBucket(bucket)
	if err != nil {
		return false, err
	}
	objectKey = normalizeObjectKey(objectKey)
	if objectKey == "" {
		return false, UploadFailed("object key is required")
	}

	_, err = client.StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}

	var respErr minio.ErrorResponse
	if errors.As(err, &respErr) && (respErr.Code == "NoSuchKey" || respErr.Code == "NoSuchObject" || respErr.Code == "NotFound") {
		return false, nil
	}

	return false, UploadFailed(fmt.Sprintf("failed to stat object %s: %v", objectKey, err))
}

func (s *MinIOStorage) GetObjectURL(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	client, err := s.ensureClient()
	if err != nil {
		return "", err
	}

	bucket, err = s.resolveBucket(bucket)
	if err != nil {
		return "", err
	}
	objectKey = normalizeObjectKey(objectKey)
	if objectKey == "" {
		return "", UploadFailed("object key is required")
	}

	url, err := client.PresignedGetObject(ctx, bucket, objectKey, expiry, nil)
	if err != nil {
		return "", UploadFailed(fmt.Sprintf("failed to generate object url for %s: %v", objectKey, err))
	}

	return url.String(), nil
}

func (s *MinIOStorage) ensureClient() (*minio.Client, error) {
	if s.client == nil {
		return nil, StorageUnavailable("minio client is not initialized")
	}

	return s.client, nil
}

func (s *MinIOStorage) resolveBucket(bucket string) (string, error) {
	bucket = normalizeObjectKey(bucket)
	if bucket != "" {
		return bucket, nil
	}

	bucket = normalizeObjectKey(s.bucket)
	if bucket == "" {
		return "", BucketUnavailable("")
	}

	return bucket, nil
}
