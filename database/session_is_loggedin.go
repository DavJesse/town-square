package database

import (
	"fmt"
	"net/http"
	"time"

	"forum/models"
)

// Checks if a user is logged in
func IsLoggedIn(r *http.Request) (*models.SessionWithUsername, bool) {
	// Try to get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// If the cookie is not found, return nil and false indicating no session
		if err == http.ErrNoCookie {
			// fmt.Println("Session cookie not present")
			return nil, false
		}
		// Handle any other error that may occur
		fmt.Println("ERROR: ", err.Error())
		return nil, false
	}

	// Query the database to retrieve the session and user information
	query := `
	SELECT s.session_id, s.user_id, s.expiry, u.username FROM
	sessions s
	INNER JOIN users u ON u.session_id = s.session_id
	WHERE s.session_id = ?`

	var session models.SessionWithUsername
	err = db.QueryRow(query, cookie.Value).Scan(&session.SessionID, &session.UserID, &session.Expiry, &session.Username)
	if err != nil {
		// If the session isn't found in the database, return nil and false
		fmt.Println("ERROR: ", err.Error())
		return nil, false
	}

	// Check if the session has expired
	if session.Expiry.Before(time.Now()) {
		fmt.Println("Session expired")
		return nil, false
	}

	// If session is found and valid, return the session object and true indicating the user is logged in
	return &session, true
}
