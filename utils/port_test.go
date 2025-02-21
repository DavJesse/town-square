package utils

import (
	"os"
	"testing"
)

/*
* tests the Port function which is responsible for getting the port number from the envirenment variable PORT
* 1. valid PORT - must be set to a number
* 2. invalid - empty, non-integer, 0(zero)
 */
func TestPort(t *testing.T) {
	table := []struct {
		name         string
		envPort      string
		expectedPort int
	}{
		{
			name:         "Port not set in environment",
			envPort:      "",
			expectedPort: 8080, // Default port
		},
		{
			name:         "Port set to a valid number",
			envPort:      "9000",
			expectedPort: 9000,
		},
		{
			name:         "Port set to an invalid number",
			envPort:      "abc", // Invalid input
			expectedPort: 8080,  // Default to 8080
		},
		{
			name:         "Port set to zero",
			envPort:      "0", // Zero is a valid integer but usually, it's an invalid port.
			expectedPort: 0,   // It's a valid case but we'll return it as is in this case.
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			// Set the environment variable for the test case
			os.Setenv("PORT", entry.envPort)
			defer os.Unsetenv("PORT") // Clean up after the test

			// Call the Port function and check the result
			result := Port()
			if result != entry.expectedPort {
				t.Errorf("For environment variable PORT='%s', expected %d but got %d", entry.envPort, entry.expectedPort, result)
			}
		})
	}
}
