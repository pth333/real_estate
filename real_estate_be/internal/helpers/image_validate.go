package helpers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"
)

const (
	maxImageSize   int64 = 10 * 1024 * 1024
	minImageWidth        = 1400
	minImageHeight       = 1050
	ratioTolerance       = 0.02
)

func ValidateImage(filename string, contentType string, fileSize int64, payload []byte) error {
	if filename == "" {
		return fmt.Errorf("filename không được để trống")
	}
	if strings.Contains(filename, "..") {
		return fmt.Errorf("filename không hợp lệ")
	}
	allowedExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	}
	if !allowedExtensions[strings.ToLower(filepath.Ext(filename))] {
		return fmt.Errorf("extension ảnh không được hỗ trợ: %s", filepath.Ext(filename))
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	normalizedContentType := strings.ToLower(strings.TrimSpace(contentType))
	if !allowedTypes[normalizedContentType] {
		return fmt.Errorf("định dạng không hợp lệ, chỉ hỗ trợ PNG, JPG, JPEG, GIF")
	}

	if fileSize <= 0 || int64(len(payload)) != fileSize {
		return fmt.Errorf("dung lượng ảnh không hợp lệ")
	}

	if fileSize > maxImageSize {
		return fmt.Errorf("dung lượng ảnh không được vượt quá 10MB")
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("file không phải ảnh hợp lệ")
	}

	if (normalizedContentType == "image/jpeg" || normalizedContentType == "image/jpg") && format != "jpeg" {
		return fmt.Errorf("nội dung file không khớp với định dạng JPEG")
	}
	if normalizedContentType == "image/png" && format != "png" {
		return fmt.Errorf("nội dung file không khớp với định dạng PNG")
	}
	if normalizedContentType == "image/gif" && format != "gif" {
		return fmt.Errorf("nội dung file không khớp với định dạng GIF")
	}

	if config.Width < minImageWidth || config.Height < minImageHeight {
		return fmt.Errorf(
			"kích thước ảnh tối thiểu là 400x300px, ảnh hiện tại %dx%dpx",
			config.Width,
			config.Height,
		)
	}

	imageRatio := float64(config.Width) / float64(config.Height)
	isSquare := config.Width == config.Height
	isFourByThree := abs(imageRatio-(4.0/3.0)) <= ratioTolerance
	if !isSquare && !isFourByThree {
		return fmt.Errorf(
			"tỷ lệ ảnh chỉ được là 4:3 hoặc hình vuông, ảnh hiện tại %dx%dpx",
			config.Width,
			config.Height,
		)
	}

	return nil
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
