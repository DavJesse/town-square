package database

import "fmt"

// UpdateUserSession updates the user's session ID in the users table.
func UpdateUserSession(sessionID string, userID int) error {
	query := `UPDATE users SET session_id = ? WHERE id = ?`
	_, err := db.Exec(query, sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to update user session ID: %w", err)
	}
	return nil
}
