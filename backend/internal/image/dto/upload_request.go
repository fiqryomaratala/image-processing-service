package dto

import "io"

type UploadRequest struct {
	OriginalFilename string
	StoredFilename   string
	ObjectKey        string
	BucketName       string
	ContentType      string
	FileSize         int64
	Width            int
	Height           int
	File             io.Reader
}
