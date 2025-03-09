package auth

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
	utils "forum/utils"
)

// LoginHandler handles user login and session creation, as well as preventing login when already logged in.
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./web/templates/index.html")
		return
	}

	// Declare utility variables
	var user models.User

	// Catch non-Get and non-POST requests
	if r.Method != "POST" {
		log.Println("METHOD ERROR: method not allowed")
		errors.MethodNotAllowedHandler(w)
		return
	}

	// Parse form data
	r.ParseMultipartForm(10 << 20)

	// Populate user credentials
	emailUsername := html.EscapeString(r.FormValue("email_username"))
	user.Password = html.EscapeString(r.FormValue("password")) // Populate password field

	// Check for empty user input
	if emailUsername == "" || user.Password == "" {
		log.Println("ERROR: Empty username or password field")
		ParseAlertMessage(w, "email and password are required")
		return
	}

	// Determine whether input is a valid email
	if utils.ValidEmail(emailUsername) {
		user.Email = emailUsername
	} else {
		user.Username = emailUsername
	}

	// Attempt to log in the user
	sessionID, err := database.LoginUser(user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("LOGIN ERROR: %v", err)
		ParseAlertMessage(w, "invalid username or password")
		return
	}

	// Login successful, set the session ID as a cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // Session expires in 24 hours
		HttpOnly: true,                           // Prevent client-side script access
		Secure:   true,                           // Ensure cookie is only sent over HTTPS
		SameSite: http.SameSiteStrictMode,        // Prevent CSRF attacks
	}

	http.SetCookie(w, &cookie)

	// Redirect the user to the home page or a protected route
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ParseAlertMessage is used for displaying alert messages in templates.
func ParseAlertMessage(w http.ResponseWriter, message string) {
	var alert models.FormError

	// Set relevant headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Parse error message
	alert.ErrorMessage = message

	if err := json.NewEncoder(w).Encode(alert); err != nil {
		errors.InternalServerErrorHandler(w)
		log.Println("JSON ENCODING ERROR: ", err)
	}
}
