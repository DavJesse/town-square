package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/models"
	"html/template"
	"strconv"
)

// GetCategoriesHandler handles requests to retrieve all categories.
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	categories, err := database.FetchCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Sends all categories
func CategoriesPageHandler(w http.ResponseWriter, r *http.Request ){
	categories, err := database.FetchCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	// Load the HTML template
	tmpl := template.Must(template.ParseFiles("web/templates/categories.html"))
	// Execute the template with the posts data
	data := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}
	tmpl.Execute(w, data)
}

// Sends all posts of a single category
func SingeCategoryPosts(w http.ResponseWriter, r *http.Request ){
	categoryID:= r.URL.Query().Get("ID")
	ID,err:=strconv.Atoi(categoryID)
	if err!=nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	posts, err := database.FetchCategoryPostsWithID(ID)
	if err != nil {
		http.Error(w, "Category doesn't exist", http.StatusInternalServerError)
		return
	}
	// Load the HTML template
	tmpl := template.Must(template.ParseFiles("web/templates/category.html"))
	// Execute the template with the posts
	data := struct {
		Posts []models.Post
	}{
		Posts: posts,
	}
	tmpl.Execute(w, data)
}