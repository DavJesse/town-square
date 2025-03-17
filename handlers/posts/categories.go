package posts

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	errors "forum/handlers/errors"

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

// Sends all posts of a single category in JSON format
func SingleCategoryPosts(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	var err error
	session, loggedIn := database.IsLoggedIn(r)

	// Retrieve user data
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			errors.InternalServerErrorHandler(w)
			log.Printf("Error getting user: %v\n", err)
			return
		}
	}

	// Extract the category ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
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
		errors.NotFoundHandler(w)
		return
	}

	// Prepare data for JSON response
	data := struct {
		Posts    []models.PostWithCategories
		IsLogged bool
		ProfPic  string
	}{
		Posts:    posts,
		IsLogged: loggedIn,
		ProfPic:  userData.Image,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
