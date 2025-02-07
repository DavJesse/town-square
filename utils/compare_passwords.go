package utils

import "golang.org/x/crypto/bcrypt"

// import "golang.org/x/crypto/bcrypt"

// Uses bcrypt to compare hashed and non-hashed passwords
func MatchPasswords(hashedPassword, inputPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		// Confirm if error is due to mismatching passwords
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
