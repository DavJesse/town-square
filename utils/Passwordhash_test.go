package utils

import (
	"forum/models"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

/*
* Tests the Passwordhash function
* It checks if the function correctly hashes a given user's password.
* 1. Valid password: Ensure that a valid password is hashed and doesn't remain the same.
* 2. Empty password: Ensure that the function can handle an empty password.
 */
func TestPasswordhash(t *testing.T) {
	table := []struct {
		name            string
		inputPassword   string
		expectedSuccess bool
	}{
		{
			name:            "Valid password",
			inputPassword:   "testPassword123",
			expectedSuccess: true,
		},
		{
			name:            "Empty password",
			inputPassword:   "",
			expectedSuccess: true,
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			user := &models.User{Password: entry.inputPassword}

			Passwordhash(user)

			// Check if the password field is correctly hashed
			if entry.expectedSuccess {
				// The password should now be hashed (not equal to the original plain password)
				if user.Password == entry.inputPassword {
					t.Errorf("expected password to be hashed, got plain password")
				}

				// Check if the hashed password can be compared with bcrypt
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(entry.inputPassword))
				if err != nil {
					t.Errorf("expected valid hash, but got error: %v", err)
				}
			} else {
				// For invalid cases, we expect the password to remain unchanged (if any such case is added in the future)
				if user.Password == entry.inputPassword {
					t.Errorf("expected password to be hashed, got plain password")
				}
			}
		})
	}
}
