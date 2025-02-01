package handlers

import (
	"forum/models"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Restrict requests to "GET" or "POST"
	if !(r.Method == "POST" || r.Method == "GET") {
		BadRequestHandler(w)
		log.Println("LoginHandler ERROR: Bad request")
		return
	}

	// Render the login template with the request object as data
	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "Failed to load Login template", http.StatusInternalServerError)
		return
	}

	// Populate user credentials
	// Determine whether input is a valid email
	if IsValidEmail(r.FormValue("email_username")) {
		user.Email = r.FormValue("email_username")
	} else {
		user.Username = r.FormValue("username")
	}

	user.Password = r.FormValue("password") // Populate password field
	log.Println("Username: ", user.Username)
	log.Println("Email: ", user.Email)
	log.Println("Password: ", user.Password)

	if err := tmpl.Execute(w, r); err != nil {
		InternalServerErrorHandler(w)
		return
	}
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.HasSuffix(email, ".com")
}
