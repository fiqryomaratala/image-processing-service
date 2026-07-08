package validator

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

func SanitizeFilename(filename string) (string, error) {
	trimmed := strings.TrimSpace(filename)
	if trimmed == "" {
		return "", FilenameInvalid()
	}

	if strings.Contains(trimmed, "/") || strings.Contains(trimmed, "\\") {
		return "", FilenameInvalid()
	}

	if filepath.Base(trimmed) != trimmed {
		return "", FilenameInvalid()
	}

	ext := filepath.Ext(trimmed)
	name := strings.TrimSuffix(trimmed, ext)
	name = sanitizeFilenamePart(name)
	if name == "" {
		return "", FilenameInvalid()
	}

	sanitizedExt := strings.ToLower(strings.TrimPrefix(ext, "."))
	if sanitizedExt == "" {
		return name, nil
	}

	return fmt.Sprintf("%s.%s", name, sanitizedExt), nil
}

func GenerateObjectKey(prefix, sanitizedFilename string) (string, error) {
	safePrefix := sanitizeObjectKeyPrefix(prefix)
	if safePrefix == "" {
		return "", FilenameInvalid()
	}

	ext := normalizeExtension(filepath.Ext(sanitizedFilename))
	if ext == "" {
		return "", InvalidExtension()
	}

	return fmt.Sprintf("%s/%s%s", safePrefix, uuid.NewString(), ext), nil
}

func sanitizeFilenamePart(value string) string {
	var builder strings.Builder
	lastUnderscore := false

	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(r)
			lastUnderscore = false
		case r == '.', r == '-', r == '_':
			builder.WriteRune(r)
			lastUnderscore = false
		case unicode.IsSpace(r):
			if !lastUnderscore {
				builder.WriteByte('_')
				lastUnderscore = true
			}
		}
	}

	return strings.Trim(builder.String(), "._-")
}

func sanitizeObjectKeyPrefix(prefix string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(prefix, "\\", "/"))
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return ""
	}

	parts := strings.Split(trimmed, "/")
	safeParts := make([]string, 0, len(parts))
	for _, part := range parts {
		safe := sanitizeFilenamePart(part)
		if safe == "" {
			continue
		}

		safeParts = append(safeParts, safe)
	}

	return strings.Join(safeParts, "/")
}
