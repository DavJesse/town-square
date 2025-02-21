package database

import (
	"fmt"
	"time"
)

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
