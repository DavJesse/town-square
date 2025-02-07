package middlewares

import (
	"context"
	"net/http"
	"time"

	"forum/database"
)

// AuthMiddleware checks for a valid session ID in the cookie before allowing access to protected routes.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session ID from the cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No session ID cookie found, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value

		// Check if the session exists and is valid
		session, err := database.GetSession(sessionID)
		if err != nil {
			// Session not found or invalid, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check if the session has expired
		if time.Now().After(session.Expiry) {
			// Session expired, delete it and redirect to login
			_ = database.DeleteSession(sessionID) // Ignore errors for now
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Attach the session to the request context for use in handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "session", session)
		r = r.WithContext(ctx)

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
