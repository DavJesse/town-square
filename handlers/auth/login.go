package auth

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
	utils "forum/utils"
)

// LoginHandler handles user login and session creation, as well as preventing login when already logged in.
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Build path to login.html
	templatePath, err := utils.GetTemplatePath("login.html")
	if err != nil {
		log.Printf("TEMPLATE AVAILABILITY ERROR: %v", err)
		errors.NotFoundHandler(w)
		return
	}

	// Parse html template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("TEMPLATE PARSING ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	// If the method is GET, if serve blank login form
	if r.Method == "GET" {
		ExecuteTemplate(w, tmpl)
		return
	}

	// Catch non-Get and non-POST requests
	if r.Method != "POST" {
		log.Println("METHOD ERROR: method not allowed")
		errors.MethodNotAllowedHandler(w)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("FORM ERROR: %v", err)
		errors.BadRequestHandler(w)
		return
	}

	// Populate user credentials
	// Determine whether input is a valid email
	emailUsername := html.EscapeString(r.FormValue("email_username"))

	if utils.ValidEmail(emailUsername) {
		user.Email = emailUsername
	} else {
		user.Username = emailUsername
	}

	// Extract form data
	user.Password = html.EscapeString(r.FormValue("password")) // Populate password field

	// Attempt to log in the user
	sessionID, err := database.LoginUser(user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("LOGIN ERROR: %v", err)
		ParseAlertMessage(w, tmpl, "invalid username or password")
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
func ParseAlertMessage(w http.ResponseWriter, tmpl *template.Template, message string) {
	// Define template path and error message
	alert := map[string]string{"ErrorMessage": message}

	// Execute the page
	err := tmpl.Execute(w, alert)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("TEMPLATE EXECUTION ERROR: %v", err)
		return
	}
}

func ExecuteTemplate(w http.ResponseWriter, tmpl *template.Template) {
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("TEMPLATE EXECUTION ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
	}
}
