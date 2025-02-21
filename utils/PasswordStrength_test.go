package utils

import (
	"testing"
)

/*
* Tests the PasswordStrength function
* It checks if the password meets the defined strength requirements:
* 1. Password length: Should be between 8 and 72 characters.
* 2. Password content: Should have at least one uppercase letter, one lowercase letter, one number, and one special character.
 */
func TestPasswordStrength(t *testing.T) {
	table := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Valid password",
			password: "Valid1Password!",
			expected: false,
		},
		{
			name:     "Password too short",
			password: "Short1!",
			expected: true,
		},
		{
			name:     "Missing uppercase letter",
			password: "validpassword1!",
			expected: true,
		},
		{
			name:     "Missing lowercase letter",
			password: "VALIDPASSWORD1!",
			expected: true,
		},
		{
			name:     "Missing digit",
			password: "ValidPassword!",
			expected: true,
		},
		{
			name:     "Missing special character",
			password: "ValidPassword1",
			expected: true,
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			err := PasswordStrength(entry.password)

			// Check if error status matches expectation
			if (err != nil) != entry.expected {
				t.Errorf("For password '%s', expected error status %v but got %v", entry.password, entry.expected, err != nil)
			}
		})
	}
}
