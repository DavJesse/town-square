package errors

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"forum/models"
	"forum/utils"
)

func TestBadRequestHandler_StatusCode400(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequestHandler(w)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBadRequestHandler_ParseTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequestHandler(w)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedContentType := "text/html; charset=utf-8"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Bad Request!") {
		t.Errorf("Expected response body to contain 'Bad Request!', but it doesn't")
	}
}

func TestBadRequestHandler_SetCodeTo400(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequestHandler(w)

	if hitch.Code != http.StatusBadRequest {
		t.Errorf("Expected hitch.Code to be %d, but got %d", http.StatusBadRequest, hitch.Code)
	}
}

func TestBadRequestHandler_SetIssueToBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequestHandler(w)

	if hitch.Issue != "Bad Request!" {
		t.Errorf("Expected hitch.Issue to be 'Bad Request!', but got '%s'", hitch.Issue)
	}
}

func TestBadRequestHandler_ParseTemplateWithHitchData(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequestHandler(w)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedHitch := models.WebError{
		Code:  http.StatusBadRequest,
		Issue: "Bad Request!",
	}

	body := w.Body.String()
	if !strings.Contains(body, strconv.Itoa(expectedHitch.Code)) || !strings.Contains(body, expectedHitch.Issue) {
		t.Errorf("Expected response body to contain %d and %s", expectedHitch.Code, expectedHitch.Issue)
	}
}

func TestNotFoundHandler_StatusCode404(t *testing.T) {
	w := httptest.NewRecorder()
	NotFoundHandler(w)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestFoundHandler_ParseTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	NotFoundHandler(w)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	expectedContentType := "text/html; charset=utf-8"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Not Found!") {
		t.Errorf("Expected response body to contain 'Not Found!', but it doesn't")
	}
}

func TestFoundHandler_SetCodeTo404(t *testing.T) {
	w := httptest.NewRecorder()
	NotFoundHandler(w)

	if hitch.Code != http.StatusNotFound {
		t.Errorf("Expected hitch.Code to be %d, but got %d", http.StatusNotFound, hitch.Code)
	}
}

func TestFoundHandler_SetIssueToNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	NotFoundHandler(w)

	if hitch.Issue != "Not Found!" {
		t.Errorf("Expected hitch.Issue to be 'Not Found!', but got '%s'", hitch.Issue)
	}
}

func TestFoundHandler_ParseTemplateWithHitchData(t *testing.T) {
	w := httptest.NewRecorder()
	NotFoundHandler(w)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	expectedHitch := models.WebError{
		Code:  http.StatusNotFound,
		Issue: "Not Found!",
	}

	body := w.Body.String()
	if !strings.Contains(body, strconv.Itoa(expectedHitch.Code)) || !strings.Contains(body, expectedHitch.Issue) {
		t.Errorf("Expected response body to contain %d and %s", expectedHitch.Code, expectedHitch.Issue)
	}
}

func TestInternalServerErrorHandler_StatusCode500(t *testing.T) {
	w := httptest.NewRecorder()
	InternalServerErrorHandler(w)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestInternalServerErrorHandler_ParseTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	InternalServerErrorHandler(w)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedContentType := "text/html; charset=utf-8"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Internal Server Error!") {
		t.Errorf("Expected response body to contain 'Internal Server Error!', but it doesn't")
	}
}

func TestInternalServerErrorHandler_SetCodeTo500(t *testing.T) {
	w := httptest.NewRecorder()
	InternalServerErrorHandler(w)

	if hitch.Code != http.StatusInternalServerError {
		t.Errorf("Expected hitch.Code to be %d, but got %d", http.StatusInternalServerError, hitch.Code)
	}
}

func TestInternalServerErrorHandler_SetIssueToNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	InternalServerErrorHandler(w)

	if hitch.Issue != "Internal Server Error!" {
		t.Errorf("Expected hitch.Issue to be 'Internal Server Error!', but got '%s'", hitch.Issue)
	}
}

func TestInternalServerErrorHandler_ParseTemplateWithHitchData(t *testing.T) {
	w := httptest.NewRecorder()
	InternalServerErrorHandler(w)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedHitch := models.WebError{
		Code:  http.StatusInternalServerError,
		Issue: "Internal Server Error!",
	}

	body := w.Body.String()
	if !strings.Contains(body, strconv.Itoa(expectedHitch.Code)) || !strings.Contains(body, expectedHitch.Issue) {
		t.Errorf("Expected response body to contain %d and %s", expectedHitch.Code, expectedHitch.Issue)
	}
}

func TestGetTemplatePath_FoundInCurrentDirectory(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the template file
	templateDir := filepath.Join(tempDir, "web", "templates")
	err = os.MkdirAll(templateDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	templateFile := "test.html"
	_, err = os.Create(filepath.Join(templateDir, templateFile))
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Change the current working directory
	oldWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// Call the function
	path, err := utils.GetTemplatePath(templateFile)
	// Check the result
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	expectedPath := filepath.Join(tempDir, "web", "templates", templateFile)
	if path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, path)
	}
}

func TestGetTemplatePath_NotFound(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change the working directory to the temp directory
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Call GetTemplatePath with a non-existent template file
	_, err = utils.GetTemplatePath("non_existent_template.html")

	// Check if an error is returned
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Check if the error message is correct
	expectedErrMsg := "template file not found: non_existent_template.html"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestGetTemplatePath_SpecialCharacters(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock project structure
	webDir := filepath.Join(tempDir, "web")
	templatesDir := filepath.Join(webDir, "templates")
	err = os.MkdirAll(templatesDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Create a template file with special characters
	specialFileName := "test@#$%^&*.html"
	specialFilePath := filepath.Join(templatesDir, specialFileName)
	_, err = os.Create(specialFilePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change working directory to the temp directory
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Test GetTemplatePath with the special character file name
	result, err := utils.GetTemplatePath(specialFileName)
	if err != nil {
		t.Errorf("GetTemplatePath returned an error: %v", err)
	}

	expected := filepath.Join(tempDir, "web", "templates", specialFileName)
	if result != expected {
		t.Errorf("Expected path %s, but got %s", expected, result)
	}
}

func TestGetTemplatePath_EmptyFileName(t *testing.T) {
	_, err := utils.GetTemplatePath("")
	if err == nil {
		t.Error("Expected an error for empty file name, but got nil")
	}
	expectedErrMsg := "template file name cannot be empty"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestGetTemplatePath_DirectoryNotExist(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change the current working directory to the temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Call GetTemplatePath
	_, err = utils.GetTemplatePath("test.html")

	// Check if the function returns an error
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Check if the error message is correct
	expectedErrMsg := "template file not found: test.html"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestGetTemplatePath_SymbolicLinks(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a symbolic link
	linkDir := filepath.Join(tempDir, "link")
	targetDir := filepath.Join(tempDir, "target")
	os.Mkdir(targetDir, 0o755)
	os.Symlink(targetDir, linkDir)

	// Create the template file
	templateDir := filepath.Join(targetDir, "web", "templates")
	os.MkdirAll(templateDir, 0o755)
	templateFile := filepath.Join(templateDir, "test.html")
	os.WriteFile(templateFile, []byte("test"), 0o644)

	// Change working directory to the symlink
	oldWd, _ := os.Getwd()
	os.Chdir(linkDir)
	defer os.Chdir(oldWd)

	// Test GetTemplatePath
	path, err := utils.GetTemplatePath("test.html")
	if err != nil {
		t.Errorf("GetTemplatePath failed: %v", err)
	}

	expectedPath := filepath.Join(targetDir, "web", "templates", "test.html")
	if path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, path)
	}
}

func TestGetTemplatePath_MultipleMatches(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple directories with the same template file
	dirs := []string{
		filepath.Join(tempDir, "project1", "web", "templates"),
		filepath.Join(tempDir, "project2", "web", "templates"),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(filepath.Join(dir, "test.html"), []byte("test content"), 0o644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Change the working directory to the deepest subdirectory
	err = os.Chdir(filepath.Join(tempDir, "project2", "web", "templates"))
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Call the function
	result, err := utils.GetTemplatePath("test.html")
	// Check the result
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := filepath.Join(tempDir, "project2", "web", "templates", "test.html")
	if result != expected {
		t.Errorf("Expected path %s, but got %s", expected, result)
	}
}
