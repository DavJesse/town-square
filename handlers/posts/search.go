package posts

import (
	"forum/database"
	// "forum/models"
	"html/template"
	"log"
	"net/http"
)

// SearchHandler handles search requests
func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter is missing", http.StatusBadRequest)
		return
	}

	posts, err := database.SearchPosts(query)
	if err != nil {
		log.Println("Error searching posts:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(posts)
	tmpl, err := template.ParseFiles("web/templates/search_results.html")
	if err != nil {
		log.Println("Failed to load template", err)
		// http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, posts); err != nil {
		log.Printf("Error executing template: %v", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
