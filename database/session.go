package database

import (
	"database/sql"
	"fmt"
	"time"

	"forum/models"
	"forum/utils"
)

// GetSession retrieves a session by its session ID.
func GetSession(sessionID string) (models.Session, error) {
	query := `SELECT session_id, user_id, expiry FROM sessions WHERE session_id = ?`
	var session models.Session
	err := db.QueryRow(query, sessionID).Scan(&session.SessionID, &session.UserID, &session.Expiry)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

// DeleteSession deletes a session by its session ID.
func DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE session_id = ?`
	_, err := db.Exec(query, sessionID)
	return err
}

// CreateSession creates a new session and updates the users table with the session ID.
func CreateSession(sessionID string, userID int, expiry time.Time) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert the new session into the sessions table
	query := `INSERT INTO sessions (session_id, user_id, expiry) VALUES (?, ?, ?)`
	_, err = tx.Exec(query, sessionID, userID, expiry)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Update the users table with the new session ID
	query = `UPDATE users SET session_id = ? WHERE id = ?`
	_, err = tx.Exec(query, sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to update user session ID: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

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

// LogoutUser handles user logout by deleting the session and clearing the session ID in the users table.
func LogoutUser(sessionID string) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete the session from the sessions table
	query := `DELETE FROM sessions WHERE session_id = ?`
	_, err = tx.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Clear the session ID in the users table
	query = `UPDATE users SET session_id = NULL WHERE session_id = ?`
	_, err = tx.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to clear user session ID: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
