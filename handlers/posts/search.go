package posts

import (
	"html/template"
	"log"
	"net/http"

	"forum/database"
	errors "forum/handlers/errors"
	"forum/models"
)

// SearchHandler handles search requests
func Search(w http.ResponseWriter, r *http.Request) {
	var user models.User

	session, logged := database.IsLoggedIn(r)

	if logged {
		user, _ = database.GetUserbySessionID(session.SessionID)
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		errors.BadRequestHandler(w)
		log.Println("Query parameter is missing")
		return
	}

	posts, err := database.SearchPosts(query)
	if err != nil {
		log.Println("Error searching posts:", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	data := struct {
		Posts    []models.Post
		IsLogged bool
		ProfPic  string
	}{
		Posts:    posts,
		IsLogged: logged,
		ProfPic:  user.Image,
	}

	// log.Println(posts) // debug log
	tmpl, err := template.ParseFiles("web/templates/search_results.html")
	if err != nil {
		log.Println("Failed to load template", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		errors.InternalServerErrorHandler(w)
	}
}
