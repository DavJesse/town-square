package database

// DeleteSession deletes a session by its session ID.
func DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE session_id = ?`
	_, err := db.Exec(query, sessionID)
	return err
}
