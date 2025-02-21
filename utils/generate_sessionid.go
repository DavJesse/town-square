package utils

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// Helper function to generate a session ID (e.g., UUID)
func GenerateSessionID() (string, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating UUID: %w", err)
	}
	return newUUID.String(), nil
}
