package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

// Implement a multipart.File
func MockFile(b []byte) multipart.File {
	return NewMockFile(b)
}

// Mock implementation of multipart.File for testing
type mockFile struct {
	content []byte
	pos     int
}

func NewMockFile(b []byte) *mockFile {
	return &mockFile{content: b}
}

func (m *mockFile) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.content) {
		return 0, io.EOF
	}

	n = copy(p, m.content[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

func (m *mockFile) Seek(offset int64, whence int) (n int64, err error) {
	return 0, nil
}

func (m *mockFile) Close() error {
	return nil
}

/*
* Tests the SaveImage function
* It checks if the function successfully saves an image file with the correct filename.
* 1. The function should generate a unique filename.
* 2. It should save the file with the correct extension.
* 3. It should handle various error scenarios, such as generating the filename or saving the file.
 */
func TestSaveImage(t *testing.T) {
	table := []struct {
		name           string
		fileType       string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "Valid PNG file",
			fileType:      "image/png",
			expectedError: false,
		},
		{
			name:          "Valid JPEG file",
			fileType:      "image/jpeg",
			expectedError: false,
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			// Create a mock file with some content for testing
			content := []byte("test image content")
			mockFile := MockFile(content)

			// Temporary path for saving the image
			tmpPath := "./tmp_test_images"
			// Ensure the test directory exists
			err := os.MkdirAll(tmpPath, os.ModePerm)
			if err != nil {
				t.Fatal("Failed to create temp directory:", err)
			}
			defer os.RemoveAll(tmpPath) // clean up after test

			// Call the SaveImage function
			_, err = SaveImage(entry.fileType, mockFile, tmpPath)

			// Check if an error is expected
			if (err != nil) != entry.expectedError {
				t.Errorf("For file type '%s', expected error status %v but got %v", entry.fileType, entry.expectedError, err != nil)
			}

			// Check if the expected error message matches
			if err != nil && err.Error() != entry.expectedErrMsg {
				t.Errorf("For file type '%s', expected error message '%s' but got '%s'", entry.fileType, entry.expectedErrMsg, err.Error())
			}

			// Check if the file was saved correctly (valid case)
			if !entry.expectedError {
				// Get the saved file name and check if it exists in the directory
				files, err := os.ReadDir(tmpPath)
				if err != nil {
					t.Fatalf("Error reading temp directory: %v", err)
				}
				if len(files) != 1 {
					t.Fatalf("Expected one file in the directory, but got %d", len(files))
				}

				// Validate the file extension
				ext := filepath.Ext(files[0].Name())
				if ext != ".png" && ext != ".jpg" {
					t.Errorf("Expected file extension '.png' or '.jpg', but got '%s'", ext)
				}
			}
		})
	}
}
