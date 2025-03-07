package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"forum/database"
	"forum/models"
)

// AuthMiddleware checks if the user has a valid session and attaches the user information to the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If we have a valid session, fetch session and user info

		session, loggedIn := database.IsLoggedIn(r)
		if !loggedIn {
			w.Header().Set("Content-Type", "application/json")
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
			r := models.Response{
				Code:     http.StatusTemporaryRedirect,
				Message:  "User is not logged in. Please log in to try again",
				Redirect: "/login",
			}
			json.NewEncoder(w).Encode(r)
			return
		}

		// Check if session is expired
		if time.Now().After(session.Expiry) {
			w.Header().Set("Content-Type", "application/json")
			// Session expired, delete it and return Unauthorized
			_ = database.DeleteSession(session.SessionID)
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			r := models.Response{
				Code:     http.StatusUnauthorized,
				Message:  "Session expired. Please log in to try again",
				Redirect: "/login",
			}
			json.NewEncoder(w).Encode(r)
			return
		}

		// Attach session and user information to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, database.SESSION_KEY, session)
		r = r.WithContext(ctx)

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
