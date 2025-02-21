package posts

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	errors "forum/handlers/errors"
	"forum/utils"

	"forum/database"
	"forum/models"
)

// GetCategoriesHandler handles requests to retrieve all categories.
func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.MethodNotAllowedHandler(w)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	categories, err := database.FetchCategories()
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("DATABASE ERROR: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Sends all categories
func CategoriesPage(w http.ResponseWriter, r *http.Request) {
	categories, err := database.FetchCategories()
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("DATABASE ERROR: %v", err)
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

	tempPath, err := utils.GetTemplatePath("category.html")
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("TEMPLATE AVAILABILITY ERROR: %v", err)
		return
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
		errors.BadRequestHandler(w)
		return
	}
	// Fetch posts from the database
	posts, err := database.FetchCategoryPostsWithID(ID)
	if err != nil {
		log.Println(err)
		errors.NotFoundHandler(w)
		return
	}
	// Load the HTML template
	tmpl, err := template.ParseFiles(tempPath)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("TEMPLATE PARSING ERROR: %v", err)
		return
	}
	// Execute the template with the posts
	data := struct {
		Posts    []models.PostWithCategories
		IsLogged bool
		ProfPic  string
	}{
		Posts:    posts,
		IsLogged: loggedIn,
		ProfPic:  userData.Image,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("TEMPLATE EXECUTION ERROR: %v", err)
	}
}
