package posts

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var userData models.User

	if r.URL.String() != "/" {
		errors.NotFoundHandler(w, r)
		return
	}

	session, loggedIn := database.IsLoggedIn(r)
	// Retrieve user data if logged in
	if !loggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userData, err := database.GetUserbySessionID(session.SessionID)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
		return
	}

	// Retrieve all posts
	posts, err := database.GetAllPosts()
	if err != nil {
		posts = []models.PostWithUsername{}
		log.Printf("Error getting posts: %v\n", err)
	}

	// Retrieve all categories
	categories, err := database.FetchCategories()
	if err != nil {
		categories = []models.Category{}
		log.Printf("Error getting categories: %v\n", err)
	}

	// Retrieve liked posts
	userLikedPosts, err := database.GetLikedPostsByUser(userData.ID)
	if err != nil {
		log.Printf("Failed to retrieve liked posts %v", err)
		userLikedPosts = []models.PostWithCategories{}
	}

	// Combine post data
	postsData := struct {
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

	// Compile post data in json response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Message: "Posts data retrieved",
		Data:    postsData,
	})
}
