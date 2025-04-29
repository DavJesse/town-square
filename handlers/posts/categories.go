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
		errors.MethodNotAllowedHandler(w, r)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	categories, err := database.FetchCategories()
	if err != nil {
		errors.InternalServerErrorHandler(w, r)
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

	// Handle Unauthorized Users
	if !loggedIn {
		WriteJSON(w, http.StatusUnauthorized, models.PostResponse{
			Message: "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
	}

	// Retrieve user data
	userData, err = database.GetUserbySessionID(session.SessionID)
	if err != nil {
		log.Printf("ERROR: cannot get user %v\n", err)
		WriteJSON(w, http.StatusInternalServerError, models.PostResponse{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Extract the category ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		log.Println("ERROR: Bad Request")
		WriteJSON(w, http.StatusBadRequest, models.PostResponse{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	categoryID := pathParts[2]
	ID, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Println("ERROR: Bad Request")
		WriteJSON(w, http.StatusBadRequest, models.PostResponse{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Retrieve all categories
	categories, err := database.FetchCategories()
	if err != nil {
		categories = []models.Category{}
		log.Printf("Error getting categories: %v\n", err)
	}

	// Fetch posts from the database
	posts, err := database.FetchCategoryPostsWithID(ID)
	if err != nil {
		errors.NotFoundHandler(w, r)
		return
	}

	// Retrieve liked posts
	userLikedPosts, err := database.GetLikedPostsByUser(userData.ID)
	if err != nil {
		log.Printf("Failed to retrieve liked posts %v", err)
		userLikedPosts = []models.PostWithCategories{}
	}

	// Prepare data for JSON response
	data := struct {
		Posts          []models.PostWithUsername   `json:"posts"`
		IsLogged       bool                        `json:"is_logged"`
		User           models.User                 `json:"user"`
		UserLikedPosts []models.PostWithCategories `json:"liked_posts"`
		Categories     []models.Category           `json:"categories"`
	}{
		Posts:          posts,
		IsLogged:       loggedIn,
		User:           userData,
		UserLikedPosts: userLikedPosts,
		Categories:     categories,
	}

	WriteJSON(w, http.StatusOK, models.Response{
		Data:    data,
		Message: "Post data retrieved",
		Code:    http.StatusOK,
	})
}
