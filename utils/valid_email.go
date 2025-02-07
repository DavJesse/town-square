package utils

import "regexp"

// func IsValidEmail(email string) bool {
// 	return strings.Contains(email, "@") && strings.HasSuffix(email, ".com")
// }

// Validate email format

func ValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}
