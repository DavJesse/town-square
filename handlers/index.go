package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/models"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch posts from the database
	posts, err := database.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	categories, err := database.FetchCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	// Load the HTML template
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	// Execute the template with the posts data
	data := struct {
		Posts      []models.PostWithUsername
		Categories []models.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	tmpl.Execute(w, data)
}
