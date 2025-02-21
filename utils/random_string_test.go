package utils

import (
	"encoding/hex"
	"testing"
)

/*
* Tests the GenerateRandomName function
* It checks if the function generates a random string of 32 characters (16 bytes).
* 1. It verifies that the generated name is 32 characters long.
* 2. It ensures the function does not return an error during execution.
* 3. It verifies that the generated name is a valid hex string.
 */
func TestGenerateRandomName(t *testing.T) {
	table := []struct {
		name           string
		expectedLength int
		expectedBool   bool
	}{
		{
			name:           "Generate valid random name",
			expectedLength: 32, // 16 bytes = 32 hex characters
			expectedBool:   false,
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			randomName, err := GenerateRandomName()

			if (err != nil) != entry.expectedBool {
				t.Errorf("expected error status %v but got %v", entry.expectedBool, err != nil)
			}

			if len(randomName) != entry.expectedLength {
				t.Errorf("expected random name length of %d, but got %d", entry.expectedLength, len(randomName))
			}

			// check if the random name is a valid hex string
			_, err = hex.DecodeString(randomName)
			if err != nil {
				t.Errorf("expected random name to be a valid hex string, but got error: %v", err)
			}
		})
	}
}
