package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// CeenerateRandomName generates a random string for the image filename.
func GenerateRandomName() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 characters in hex
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
