package posts

import (
	"html/template"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/" {
		errors.NotFoundHandler(w, r)
		return
	}

	response := models.Response{
		Code: 200,
	}

	session, loggedIn := database.IsLoggedIn(r)
	// Retrieve user data if logged in
	if !loggedIn {
		response.Message = "Please login to proceed"
		response.Redirect = "/login"
	} else {
		userData, err := database.GetUserbySessionID(session.SessionID)
		if err != nil {
			log.Printf("Error getting user: %v\n", err) // Add error logging
			response.Redirect = "/login"
		} else {
			response.Data = userData
		}
	}

	// Load the HTML template
	// Parse template with function to replace '\n' with '<br>'
	tmpl := template.Must(template.ParseFiles("./web/templates/index.html"))

	tmpl.Execute(w, response)
}
