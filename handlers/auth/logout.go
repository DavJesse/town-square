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
		log.Println("METHOD ERROR: method not allowed")
		errors.MethodNotAllowedHandler(w)
		return
	}

	// Destroy cookies if any
	err := LogOutSession(w, r)
	if err != nil {
		log.Printf("LOG OUT ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	// redirect client to home
	http.Redirect(w, r, "/", http.StatusFound)
}

func LogOutSession(w http.ResponseWriter, r *http.Request) error {
	hasCookie, cookie, err := database.HasCookie(r)
	if err != nil {
		log.Println("COOKIE ERROR: ", err)
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

	} else {
		return nil
	}

	// Delete session from database
	sessionID := cookie.Value
	err = database.DeleteSession(sessionID)
	if err != nil {
		log.Println("DATABASE ERROR: ", err)
		return err
	}

	return nil
}
