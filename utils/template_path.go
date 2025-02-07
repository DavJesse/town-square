package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetTemplatePath(templateFile string) (string, error) {
	// catch empty template files
	if templateFile == "" {
		return "", fmt.Errorf("template file name cannot be empty")
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up until we find the project root
	dir := wd
	for {
		// Construct path to template, check if constructed path exists
		templatePath := filepath.Join(dir, "web", "templates", templateFile)
		if _, err := os.Stat(templatePath); err == nil {
			return templatePath, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir { // Stop if we reach the root
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("template file not found: %s", templateFile)
}
