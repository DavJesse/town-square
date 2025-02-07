package utils

import (
	"fmt"
	"regexp"
)

// PasswordStrength checks password strength and length constraints
func PasswordStrength(password string) error {
	// Ensure password does not exceed 72 bytes
	if len([]byte(password)) > 72 {
		return fmt.Errorf("your password exceeds the maximum length of 72 bytes")
	}

	// Ensure password is at least 8 characters long
	if len(password) < 8 {
		return fmt.Errorf("password is too short: it should have at least 8 characters")
	}

	// Check for at least one uppercase letter/lowercase/one digit/one special character
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

	// Validate all conditions
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return fmt.Errorf("password too weak: it should have at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}
