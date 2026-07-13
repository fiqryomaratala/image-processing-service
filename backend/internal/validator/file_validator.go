package validator

import (
	"fmt"
	"io"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"
)

type FileValidator struct {
	config config.UploadConfig
}

func NewFileValidator(cfg config.UploadConfig) *FileValidator {
	return &FileValidator{
		config: cfg,
	}
}

func (v *FileValidator) Validate(upload filepkg.Upload) (*filepkg.ValidationResult, error) {
	sanitizedFilename, err := SanitizeFilename(upload.Filename)
	if err != nil {
		return nil, err
	}

	content, detectedSize, err := readContent(upload.Reader, v.config.MaxSizeBytes)
	if err != nil {
		return nil, err
	}

	size := upload.Size
	if size <= 0 {
		size = detectedSize
	}

	if size > v.config.MaxSizeBytes {
		return nil, FileTooLarge()
	}

	extension := extensionFromFilename(sanitizedFilename)
	if !isAllowedExtension(extension, v.config.AllowedExtensions) {
		return nil, InvalidExtension()
	}

	mimeType := DetectMIMEType(content)
	if !isAllowedMIMEType(mimeType, v.config.AllowedMIMETypes) {
		return nil, InvalidMimeType()
	}

	if !extensionsMatch(mimeType, extension) {
		return nil, InvalidExtension()
	}

	objectKey, err := GenerateObjectKey(v.config.ObjectKeyPrefix, sanitizedFilename)
	if err != nil {
		return nil, err
	}

	return &filepkg.ValidationResult{
		OriginalFilename:  upload.Filename,
		SanitizedFilename: sanitizedFilename,
		ObjectKey:         objectKey,
		Extension:         extension,
		MIMEType:          mimeType,
		Size:              size,
		Content:           content,
	}, nil
}

func readContent(reader io.Reader, maxSize int64) ([]byte, int64, error) {
	if reader == nil {
		return nil, 0, InvalidImage()
	}

	limited := io.LimitReader(reader, maxSize+1)
	content, err := io.ReadAll(limited)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file content: %w", err)
	}

	if int64(len(content)) > maxSize {
		return nil, 0, FileTooLarge()
	}

	return content, int64(len(content)), nil
}
