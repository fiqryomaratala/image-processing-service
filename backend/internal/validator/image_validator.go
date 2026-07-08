package validator

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"
	_ "golang.org/x/image/webp"
)

type ImageValidator struct {
	fileValidator *FileValidator
	config        config.UploadConfig
}

func NewImageValidator(cfg config.UploadConfig) *ImageValidator {
	return &ImageValidator{
		fileValidator: NewFileValidator(cfg),
		config:        cfg,
	}
}

func (v *ImageValidator) Validate(upload filepkg.Upload) (*filepkg.ValidationResult, error) {
	result, err := v.fileValidator.Validate(upload)
	if err != nil {
		return nil, err
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(result.Content))
	if err != nil {
		return nil, InvalidImage().WithCause(err)
	}

	if err := validateImageDimensions(cfg.Width, cfg.Height, v.config); err != nil {
		return nil, err
	}

	if _, _, err := image.Decode(bytes.NewReader(result.Content)); err != nil {
		return nil, InvalidImage().WithCause(err)
	}

	result.Image = &filepkg.ImageMetadata{
		Width:  cfg.Width,
		Height: cfg.Height,
		Format: format,
	}

	return result, nil
}

func validateImageDimensions(width, height int, cfg config.UploadConfig) error {
	if width < cfg.MinWidth || height < cfg.MinHeight {
		return ImageTooLarge()
	}

	if width > cfg.MaxWidth || height > cfg.MaxHeight {
		return ImageTooLarge()
	}

	return nil
}
