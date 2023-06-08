package helper

import (
	"mime/multipart"
	"net/http"
)

func IsFileTypeAllowed(fileHeader *multipart.FileHeader) bool {
	allowedTypes := []string{"application/pdf"}
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	fileType := http.DetectContentType(buffer)
	for _, allowedType := range allowedTypes {
		if fileType == allowedType {
			return true
		}
	}

	return false
}

// Fungsi untuk memeriksa ukuran file yang diunggah
func IsFileSizeAllowed(fileHeader *multipart.FileHeader) bool {
	maxSize := int64(5 * 1024 * 1024) // 5MB
	fileSize := fileHeader.Size
	return fileSize <= maxSize
}
