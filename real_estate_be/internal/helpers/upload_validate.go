package helpers

import (
	"fmt"
	"path/filepath"
	"strings"
)

var allowedContentTypes = map[string]bool{
	"image/jpeg":       true,
	"image/png":        true,
	"image/gif":        true,
	"image/jpg":        true,
	"video/mp4":        true,
	"video/quicktime":  true,
}

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".mp4": true, ".mov": true,
}

func ValidateContentType(ct string) error {
	if !allowedContentTypes[strings.ToLower(ct)] {
		return fmt.Errorf("content_type không được hỗ trợ: %s", ct)
	}
	return nil
}

func ValidateFilename(name string) error {
	if name == "" {
		return fmt.Errorf("filename không được để trống")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("filename không hợp lệ")
	}
	ext := strings.ToLower(filepath.Ext(name))
	if !allowedExtensions[ext] {
		return fmt.Errorf("extension file không được hỗ trợ: %s", ext)
	}
	return nil
}
