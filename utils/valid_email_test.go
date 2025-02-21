package utils

import "testing"

/*
* Tests the ValidEmail function
* It checks the output of the function given various valid and invalid email forms:
* 1. Valid emails: Include lowercase, uppercase, emails with digits, and special characters like +.
* 2. Invalid emails: Include emails with missing @, multiple @, missing domain, spaces, and special characters in the domain.
 */
func TestValidEmail(t *testing.T) {
	table := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid email with lowercase",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Valid email with uppercase",
			email:    "Test@Example.com",
			expected: true,
		},
		{
			name:     "Valid email with digits",
			email:    "test123@example.com",
			expected: true,
		},
		{
			name:     "Valid email with special characters",
			email:    "test.email+alex@leetcode.com",
			expected: true,
		},
		{
			name:     "Invalid email with missing '@'",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "Invalid email with multiple '@'",
			email:    "test@@example.com",
			expected: false,
		},
		{
			name:     "Invalid email with missing domain",
			email:    "test@.com",
			expected: false,
		},
		{
			name:     "Invalid email with spaces",
			email:    "test @example.com",
			expected: false,
		},
		{
			name:     "Invalid email with special character in domain",
			email:    "test@example!com",
			expected: false,
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			result := ValidEmail(entry.email)
			if result != entry.expected {
				t.Errorf("For email '%s', expected %v but got %v", entry.email, entry.expected, result)
			}
		})
	}
}
