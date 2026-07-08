package dto

import "io"

type UploadRequest struct {
	OriginalFilename string
	ContentType      string
	FileSize         int64
	File             io.Reader
}
