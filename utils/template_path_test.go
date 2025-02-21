package utils

import (
	"os"
	"testing"
)

/*
* Tests the GetTemplatePath function.
* It checks the output for different template file scenarios:
* 1. Valid template file: When the template file is valid, it should return the correct path.
* 2. Empty template file name: If the template file name is empty, the function should return an error.
* 3. Template file not found: If the template file is not found in the expected directory structure, the function should return an error.
 */
func TestGetTemplatePath(t *testing.T) {
	workingDir := ""
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		t.Errorf("An error occured when getting the current working directory")
	} else {

		table := []struct {
			name           string
			templateFile   string
			mockCurrentDir string
			shouldError    bool
		}{
			{
				name:           "Valid template file",
				templateFile:   "index.html",
				mockCurrentDir: workingDir,
				shouldError:    false,
			},
			{
				name:           "Empty template file name",
				templateFile:   "",
				mockCurrentDir: workingDir,
				shouldError:    true,
			},
			{
				name:           "Template file not found",
				templateFile:   "missing-template.html",
				mockCurrentDir: workingDir,
				shouldError:    true,
			},
		}

		for _, entry := range table {
			t.Run(entry.name, func(t *testing.T) {
				// Mock the current directory
				os.Setenv("PWD", entry.mockCurrentDir)
				defer os.Unsetenv("PWD")

				// Call the GetTemplatePath function
				result, err := GetTemplatePath(entry.templateFile)

				// Check if an error is expected
				if (err != nil) != entry.shouldError {
					t.Errorf("For template file '%s', expected error status %v but got %v", entry.templateFile, entry.shouldError, err != nil)
				}

				// Only check the path if no error is expected
				if !entry.shouldError && result == "" {
					t.Errorf("For template file '%s', expected a valid path but got an empty string", entry.templateFile)
				}
			})
		}
	}
}
