package utils

import (
	"path/filepath"
	"strings"
)

func IsAllowedFileType(filePath string, allowedTypes []string) bool {
    ext := strings.ToLower(filepath.Ext(filePath))
    extAllowed := false
    for _, allowed := range allowedTypes {
        if ext == strings.ToLower(allowed) {
            extAllowed = true
            break
        }
    }
    if !extAllowed {
        return false
    }
    return false
}