package handlers

import (
	"log"
	"net/http"
	"text/template"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
	"forum/utils"
)

func SPA(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	var err error
	session, loggedIn := database.IsLoggedIn(r)

	// Retrieve user data if logged in
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			log.Printf("Error getting user: %v\n", err) // Add error logging
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	if r.URL.String() != "/spa" {
		errors.NotFoundHandler(w, r)
		return
	}

	// Fetch posts from the database
	posts, err := database.GetAllPosts()
	if err != nil {
		errors.InternalServerErrorHandler(w, r)
		return
	}
	categories, err := database.FetchCategories()
	if err != nil {
		errors.InternalServerErrorHandler(w, r)
		return
	}
	// Load the HTML template
	// Parse template with function to replace '\n' with '<br>'
	tmpl := template.Must(template.New("spa.html").Funcs(template.FuncMap{
		"replaceNewlines": utils.ReplaceNewlines,
	}).ParseFiles("./web/templates/spa.html"))

	// Execute the template with the posts data
	data := struct {
		Posts      []models.PostWithUsername
		Categories []models.Category
		IsLogged   bool
		ProfPic    string
	}{
		Posts:      posts,
		Categories: categories,
		IsLogged:   loggedIn,
		ProfPic:    userData.Image,
	}
	tmpl.Execute(w, data)
}
