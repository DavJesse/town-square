package posts

import (
	"log"
	"net/http"

	"forum/database"
	"forum/models"
)

func IndexRouteHandler(w http.ResponseWriter, r *http.Request) {
	// Always serve the HTML template for the root path
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./web/templates/index.html")
		return
	}
}

func GetIndexData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userData models.User
	session, loggedIn := database.IsLoggedIn(r)
	// Retrieve user data if logged in
	if !loggedIn {
		WriteJSON(w, http.StatusUnauthorized, models.PostResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	userData, err := database.GetUserbySessionID(session.SessionID)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		WriteJSON(w, http.StatusInternalServerError, models.PostResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
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

	// Get user's posts
	userPosts, err := database.PostsFilterByUser(userData.ID)
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		userPosts = []models.PostWithUsername{}
	}

	// Combine post data
	postsData := struct {
		User           models.User                 `json:"user"`
		IsLogged       bool                        `json:"is_logged"`
		Categories     []models.Category           `json:"categories"`
		AllPosts       []models.PostWithUsername   `json:"all_posts"`
		UserLikedPosts []models.PostWithCategories `json:"liked_posts"`
		UserPosts      []models.PostWithUsername   `json:"user_posts"`
	}{
		User:           userData,
		IsLogged:       loggedIn,
		Categories:     categories,
		AllPosts:       posts,
		UserLikedPosts: userLikedPosts,
		UserPosts:      userPosts,
	}

	// Compile post data in json response
	WriteJSON(w, http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "Posts data retrieved",
		Data:    postsData,
	})
}
