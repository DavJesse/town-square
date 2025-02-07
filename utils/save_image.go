package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveImage(fileType string, file multipart.File, path string) (string, error) {
	randomFileName, err := GenerateRandomName()
	if err != nil {
		return "", errors.New("failed to generate a unique filename")
	}

	// Determine the file extension based on the MIME type
	var ext string
	switch fileType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpg"
	}

	// Construct the full filename
	filename := randomFileName + ext

	filePath := filepath.Join(path, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("failed to save the file")
	}
	defer outFile.Close()

	// Copy the file content to the new file
	if _, err := io.Copy(outFile, file); err != nil {
		return "", errors.New("failed to save the file")
	}
	return filename, nil
}
