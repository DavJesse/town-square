package database

import "fmt"

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
