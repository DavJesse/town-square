package utils

import (
	"fmt"
	"github.com/gofrs/uuid"
)

// GenerateUUID generates a new version 4 UUID and returns it as a string
func GenerateUUID() (string, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating UUID: %w", err)
	}
	return newUUID.String(), nil
}
