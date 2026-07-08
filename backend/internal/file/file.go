package file

import "io"

type Upload struct {
	Filename string
	Size     int64
	Reader   io.Reader
}

type ImageMetadata struct {
	Width  int
	Height int
	Format string
}

type ValidationResult struct {
	OriginalFilename  string
	SanitizedFilename string
	ObjectKey         string
	Extension         string
	MIMEType          string
	Size              int64
	Content           []byte
	Image             *ImageMetadata
}
