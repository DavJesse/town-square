package posts

import (
	errors "forum/handlers/errors"
	utils "forum/utils"
	"html/template"
	"log"
	"net/http"

	"forum/database"
	"forum/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
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

	// Fetch posts from the database
	posts, err := database.GetAllPosts()
	if err != nil {
		errors.InternalServerErrorHandler(w)
		return
	}
	categories, err := database.FetchCategories()
	if err != nil {
		errors.InternalServerErrorHandler(w)
		return
	}
	// Load the HTML template
	// Parse template with function to replace '\n' with '<br>'
	tmpl := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"replaceNewlines": utils.ReplaceNewlines,
	}).ParseFiles("./web/templates/index.html"))

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
