package database

import (
	"forum/models"
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
