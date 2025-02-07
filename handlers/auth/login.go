package auth

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
	"forum/utils"
)

// LoginHandler handles user login and session creation, as well as preventing login when already logged in.
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Build path to login.html
	templatePath, err := utils.GetTemplatePath("login.html")
	if err != nil {
		log.Println("Not nil")
		errors.InternalServerErrorHandler(w)
		log.Println("Could not find template file: ", err)
		return
	}

	// If the method is GET, if serve blank login form
	if r.Method == "GET" {
		// Skip login for users who are already loged in
		{
			hasCookie, _, err := HasCookie(r)
			if err != nil {
				errors.InternalServerErrorHandler(w)
				log.Println("Error checking session cookie: ", err)
			}
			if hasCookie {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		// Render login form
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			http.Error(w, "Failed to load Login template", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			errors.InternalServerErrorHandler(w)
			return
		}
		return
	}

	// Catch non-Get and non-POST requests
	if r.Method != "POST" {
		errors.BadRequestHandler(w)
		log.Println("LoginHandler ERROR: Bad request method")
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Populate user credentials
	// Determine whether input is a valid email
	if utils.ValidEmail(r.FormValue("email_username")) {
		user.Email = r.FormValue("email_username")
	} else {
		user.Username = r.FormValue("email_username")
	}

	// Extract form data
	user.Password = r.FormValue("password") // Populate password field
	emailUsername := r.FormValue("email_username")
	password := r.FormValue("password")

	// Validate credentials
	_, err = database.VerifyUser(emailUsername, password)
	if err != nil {
		issue := err.Error()
		// Prepare template to render error message
		tmpl, err1 := template.ParseFiles(templatePath)
		if err1 != nil {
			errors.InternalServerErrorHandler(w)
			return
		}

		// Check for non-existent users
		if strings.Contains(issue, "user does not exist") {
			ParseAlertMessage(w, tmpl, "user does not exist")
			log.Println("INFO: User does not exist")
			return
		}

		ParseAlertMessage(w, tmpl, "wrong password")
		log.Printf("ERROR: %v", err)
		return
	}

	// Attempt to log in the user
	sessionID, err := database.LoginUser(user.Username, user.Email, user.Password)
	if err != nil {
		// Login failed, render the login template with an error message
		templatePath, err := utils.GetTemplatePath("login.html")
		if err != nil {
			errors.InternalServerErrorHandler(w)
			log.Println("Could not find template file: ", err)
			return
		}

		// Render error message
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			http.Error(w, "Failed to load Login template", http.StatusInternalServerError)
			return
		}

		// Pass the error message to the template
		data := struct {
			Error string
		}{
			Error: "Invalid username/email or password",
		}
		if err := tmpl.Execute(w, data); err != nil {
			errors.InternalServerErrorHandler(w)
			return
		}

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

func HasCookie(r *http.Request) (bool, *http.Cookie, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil, nil // No session cookie, but no error
		}
		return false, nil, err // Actual error (e.g., internal failure)
	}
	return true, cookie, nil // Cookie exists
}

// ParseAlertMessage is used for displaying alert messages in templates.
func ParseAlertMessage(w http.ResponseWriter, tmpl *template.Template, message string) {
	// Define template path and error message
	alert := map[string]string{"ErrorMessage": message}
	// Execute the page
	err := tmpl.Execute(w, alert)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("Could not execute template %v", err)
		return
	}
}

func Success(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errors.BadRequestHandler(w)
		return
	}
	http.Error(w, "Loged in1!", http.StatusInternalServerError)
}
