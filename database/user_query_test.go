package database

import (
	"fmt"
	"strings"
	"testing"
)

func TestVerifyUser_ValidCredentials(t *testing.T) {
	result, _ := VerifyUser("milton@mail.com", "mPass")

	if !result {
		t.Errorf("Expected VerifyUser to return true for valid credentials, but got false")
	}
}

func TestVerifyUser_InvalidCredentials(t *testing.T) {
	result, _ := VerifyUser("milton@mail.com", "password")

	if result {
		t.Errorf("Expected VerifyUser to return false for valid credentials, but got true")
	}
}

func TestGetUserByMailOrName(t *testing.T) {
	var err error

	user, err := GetUserByEmailOrUsername("milton@mail.com", "milton")
	if user.Username != "milton" || err != nil {
		fmt.Println("USername: ", user.Username, " Passcode: ", user.Password, " Email: ", user.Email)
		t.Errorf("expected username to be milton, but got %s\n", user.Username)
	}
}

func TestVerifyUser_NonExistentUser(t *testing.T) {
	result, _ := VerifyUser("nonexistent@mail.com", "password")

	if result {
		t.Errorf("Expected VerifyUser to return false for non-existent user, but got true")
	}
}

func TestVerifyUser_NonExistentUsername(t *testing.T) {
	result, _ := VerifyUser("nonexistent@mail.com", "password")

	if result {
		t.Errorf("Expected VerifyUser to return false for non-existent username, but got true")
	}
}

func TestVerifyUser_CaseSensitiveEmail(t *testing.T) {
	// Test with a lowercase email
	resultLower, _ := VerifyUser("milton@mail.com", "mPass")
	if !resultLower {
		t.Errorf("Expected VerifyUser to return true for lowercase email, but got false")
	}

	// Test with an uppercase email
	resultUpper, _ := VerifyUser("MILTON@MAIL.COM", "mPass")
	if resultUpper {
		t.Errorf("Expected VerifyUser to return false for uppercase email, but got true")
	}
}

func TestVerifyUser_CaseSensitiveUsername(t *testing.T) {
	// Test with correct username but different case
	result, _ := VerifyUser("MILTON@mail.com", "mPass")

	if result {
		t.Errorf("Expected VerifyUser to return false for case-sensitive username, but got true")
	}

	// Test with correct username and correct case
	result, _ = VerifyUser("milton@mail.com", "mPass")

	if !result {
		t.Errorf("Expected VerifyUser to return true for correct case-sensitive username, but got false")
	}
}

func TestVerifyUser_MaxLengthInputs(t *testing.T) {
	maxLengthEmail := "a" + strings.Repeat("b", 254) + "@c.com" // 256 characters
	maxLengthPassword := strings.Repeat("x", 72)                // 72 characters (common max length for bcrypt)

	result, _ := VerifyUser(maxLengthEmail, maxLengthPassword)

	if result {
		t.Errorf("Expected VerifyUser to return false for max length inputs, but got true")
	}
}
