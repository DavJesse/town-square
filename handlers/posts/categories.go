package posts

import (
	"encoding/json"
	errors "forum/handlers/errors"
	"log"
	"net/http"

	"forum/database"
	"forum/models"
	"html/template"
	"strconv"
	"strings"
)

// GetCategoriesHandler handles requests to retrieve all categories.
func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Invalid request method")
		errors.BadRequestHandler(w)
		return
	}

	categories, err := database.FetchCategories()
	if err != nil {
		log.Println("Failed to fetch categories")
		errors.InternalServerErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Sends all categories
func CategoriesPage(w http.ResponseWriter, r *http.Request) {
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
func SingeCategoryPosts(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	var err error
	session, loggedIn := database.IsLoggedIn(r)

	// Retrieve user data
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			log.Printf("Error getting user: %v\n", err) // Add error logging
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// Extract the category ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		log.Println("Invalid request method")
		errors.BadRequestHandler(w)
		return
	}
	categoryID := pathParts[2]
	ID, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Fetch posts from the database
	posts, err := database.FetchCategoryPostsWithID(ID)
	if err != nil {
		log.Println("Category doesn't exist")
		errors.InternalServerErrorHandler(w)
		return
	}
	// Load the HTML template
	tmpl, err := template.ParseFiles("web/templates/category.html")
	if err != nil {
		log.Println("Error: ", err)
		errors.InternalServerErrorHandler(w)
		return
	}
	// Execute the template with the posts
	data := struct {
		Posts    []models.Post
		IsLogged bool
		ProfPic  string
	}{
		Posts:    posts,
		IsLogged: loggedIn,
		ProfPic:  userData.Image,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
