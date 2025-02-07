package database

import (
	"database/sql"
	"fmt"
	"time"

	"forum/utils"
)

// LoginUser handles user login and session management.
func LoginUser(detailOne, detailTwo, password string) (string, error) {
	// Fetch the user by username or email
	query := `SELECT id, password, session_id FROM users WHERE username = ? OR username = ? OR email = ? OR email = ?`
	var userID int
	var hashedPassword string
	var existingSessionID sql.NullString
	err := db.QueryRow(query, detailOne, detailTwo, detailOne, detailTwo).Scan(&userID, &hashedPassword, &existingSessionID)
	if err != nil {
		return "", fmt.Errorf("invalid username/email or password: %w", err)
	}

	// Verify the password (assuming you have a function to compare hashed passwords)
	if match, _ := utils.MatchPasswords(hashedPassword, password); !match {
		return "", fmt.Errorf("invalid username/email or password")
	}

	// If the user already has a session, destroy it
	if existingSessionID.Valid {
		err = DeleteSession(existingSessionID.String)
		if err != nil {
			return "", fmt.Errorf("failed to delete existing session: %w", err)
		}
	}

	// Generate a new session ID
	newSessionID, err := utils.GenerateSessionID()
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Set session expiry (e.g., 24 hours from now)
	expiry := time.Now().Add(24 * time.Hour)

	// Create the new session
	err = CreateSession(newSessionID, userID, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return newSessionID, nil
}
