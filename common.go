package main

import (
	"path/filepath"
	"strings"
)

func IsSupported(path string) bool {
	ext := filepath.Ext(strings.ToLower(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".bmp"
}

func ToBmp(filePath string) string {
	ext := filepath.Ext(filePath)
	newPath := filePath[0:len(filePath)-len(ext)] + ".bmp"
	return newPath
}
