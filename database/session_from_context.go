package database

import (
	"net/http"

	"forum/models"
)

// GetSessionFromContext retrieves the session from the context using the same key
func SessionFromContext(r *http.Request) (*models.SessionWithUsername, bool) {
	// retrieve the session from the context
	session, ok := r.Context().Value(SESSION_KEY).(*models.SessionWithUsername)
	if !ok {
		// if the session is not found or cannot be cast, return false
		return nil, false
	}
	return session, true
}
