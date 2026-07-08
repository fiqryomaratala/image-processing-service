package validator

import (
	"net/http"
	"path/filepath"
	"strings"
)

var mimeExtensionMap = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func DetectMIMEType(content []byte) string {
	sniffLen := 512
	if len(content) < sniffLen {
		sniffLen = len(content)
	}

	return strings.ToLower(http.DetectContentType(content[:sniffLen]))
}

func normalizeExtension(extension string) string {
	trimmed := strings.ToLower(strings.TrimSpace(extension))
	if trimmed == "" {
		return ""
	}

	if !strings.HasPrefix(trimmed, ".") {
		return "." + trimmed
	}

	return trimmed
}

func extensionFromFilename(filename string) string {
	return normalizeExtension(filepath.Ext(filename))
}

func extensionFromMIMEType(mimeType string) string {
	return mimeExtensionMap[strings.ToLower(strings.TrimSpace(mimeType))]
}

func extensionsMatch(detectedMIMEType, filenameExtension string) bool {
	detectedExtension := extensionFromMIMEType(detectedMIMEType)
	normalizedFilenameExtension := normalizeExtension(filenameExtension)
	if detectedExtension == "" || normalizedFilenameExtension == "" {
		return false
	}

	if detectedExtension == normalizedFilenameExtension {
		return true
	}

	if detectedMIMEType == "image/jpeg" {
		return normalizedFilenameExtension == ".jpg" || normalizedFilenameExtension == ".jpeg"
	}

	return false
}

func isAllowedExtension(extension string, allowed []string) bool {
	normalized := strings.TrimPrefix(normalizeExtension(extension), ".")
	for _, item := range allowed {
		if normalized == strings.TrimPrefix(normalizeExtension(item), ".") {
			return true
		}
	}

	return false
}

func isAllowedMIMEType(mimeType string, allowed []string) bool {
	normalized := strings.ToLower(strings.TrimSpace(mimeType))
	for _, item := range allowed {
		if normalized == strings.ToLower(strings.TrimSpace(item)) {
			return true
		}
	}

	return false
}
