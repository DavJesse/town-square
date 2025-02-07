package utils

import "strings"

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.HasSuffix(email, ".com")
}
