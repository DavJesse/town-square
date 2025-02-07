package auth

import (
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/handlers/errors"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// Restrict non-POST requests
	if r.Method != "POST" {
		errors.BadRequestHandler(w)
		log.Println("ERROR: LogoutHandler does not allow non-POST requests")
		return
	}

	// Destroy cookies if any
	err := LogOutSession(w, r)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		return
	}

	// redirect client to home
	http.Redirect(w, r, "/home", http.StatusFound)
}

func LogOutSession(w http.ResponseWriter, r *http.Request) error {
	hasCookie, cookie, err := HasCookie(r)
	if err != nil {
		log.Println("ERROR: checking session cookie failed: ", err)
		return err
	}

	// Invalidate session cookie
	if hasCookie {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),         // Expires imediately
			MaxAge:   0,                       // Explicitly expire cookie
			HttpOnly: true,                    // Restricts against JavaScript access
			Secure:   true,                    // Ensures cookies only sent over HTTPS
			SameSite: http.SameSiteStrictMode, // Restricts against CSRF requests
		})

		log.Println("INFO: session cookie has been invalidate")
	} else {
		log.Println("INFO: session cookie not found")
		return nil
	}

	// Delete session from database
	sessionID := cookie.Value
	err = database.DeleteSession(sessionID)
	if err != nil {
		log.Println("ERROR: deleting session from database failed: ", err)
		return err
	}

	return nil
}
