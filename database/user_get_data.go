package database

import (
	"database/sql"
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
	cookieExists, cookie, err := Cookie(r)
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

	// Return user data
	return userData.ID, userData.Username, nil
}

func Cookie(r *http.Request) (bool, *http.Cookie, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil, nil // No session cookie, but no error
		}
		return false, nil, err // Actual error (e.g., internal failure)
	}
	return true, cookie, nil // Cookie exists
}

// HasCookie checks if the session cookie exists and matches the session_id in the users table.
func HasCookie(r *http.Request) (bool, *http.Cookie, error) {
	// Get the session_id from the cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil, nil // No session cookie, but no error
		}
		return false, nil, err // Actual error (e.g., internal failure)
	}

	// If we have a cookie, check if it matches the session_id in the users table
	var userID int
	var storedSessionID sql.NullString

	// Query the users table for the session_id corresponding to the session cookie
	query := `SELECT id, session_id FROM users WHERE session_id = ?`
	err = db.QueryRow(query, cookie.Value).Scan(&userID, &storedSessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No matching session ID in the users table
			return false, cookie, nil
		}
		return false, nil, fmt.Errorf("error checking session_id: %w", err)
	}

	// If there's no session_id in the database, or it doesn't match the cookie, return false
	if !storedSessionID.Valid || storedSessionID.String != cookie.Value {
		return false, cookie, nil
	}

	// The session_id in the cookie matches the one in the database
	return true, cookie, nil
}
