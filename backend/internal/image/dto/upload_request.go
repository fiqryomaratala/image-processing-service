package dto

import filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"

type UploadRequest struct {
	File filepkg.Upload
}
