package database

import (
	"errors"
	"fmt"
	"net/http"
)

// Assuming that UserData is a struct that holds the user data.
type UserData struct {
	ID       int
	Username string
}

// GetUserData retrieves user data from a session and stores it in the request context.
func GetUserData(r *http.Request) (int, string, error) {
	cookieExists, cookie, err := HasCookie(r)
	if err != nil {
		fmt.Println("Redirected to login")
		return 0, "", errors.New("user is not logged in")
	}

	if !cookieExists {
		return 0, "", errors.New("no such cookie")
	}

	userData, err := GetUserbySessionID(cookie.Value)
	// fmt.Printf("UserData retrieved: %+v\n", userData) // Add debug logging
	if err != nil {
		msg := fmt.Sprintf("Error getting user: %v\n", err)
		return 0, "", errors.New(msg)
	}

	// Save user data to the request context
	// ctx := r.Context()
	// ctx = context.WithValue(ctx, "userID", userData.ID)
	// ctx = context.WithValue(ctx, "username", userData.Username)

	// // Pass the updated context back to the request
	// r = r.WithContext(ctx)

	// Return user data
	return userData.ID, userData.Username, nil
}

func HasCookie(r *http.Request) (bool, *http.Cookie, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil, nil // No session cookie, but no error
		}
		return false, nil, err // Actual error (e.g., internal failure)
	}
	return true, cookie, nil // Cookie exists
}
