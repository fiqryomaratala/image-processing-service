package storage

import "strings"

func normalizeObjectKey(objectKey string) string {
	return strings.Trim(strings.TrimSpace(objectKey), "/")
}
