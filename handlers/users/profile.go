package auth

import (
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/posts"
	"forum/models"
)

func ProfileRouteHandle(w http.ResponseWriter, r *http.Request) {
	// Always serve the HTML template for the root path
	if r.URL.Path == "/profile" {
		http.ServeFile(w, r, "./web/templates/index.html")
		return
	}
}

// ViewUserProfile handler
func GetProfileData(w http.ResponseWriter, r *http.Request) {
	session, logged := database.IsLoggedIn(r)
	if !logged {
		posts.WriteJSON(w, http.StatusUnauthorized, models.PostResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	userData, err := database.GetUserbySessionID(session.SessionID)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		posts.WriteJSON(w, http.StatusInternalServerError, models.PostResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
		return
	}

	// Get user's posts
	userPosts, err := database.PostsFilterByUser(userData.ID)
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		userPosts = []models.Post{}
	}

	// Get user's liked posts
	userLikedPosts, err := database.GetLikedPostsByUser(userData.ID)
	if err != nil {
		log.Printf("Failed to retrieve liked posts %v", err)
		userLikedPosts = []models.PostWithCategories{}
	}

	// Combine user data and user posts
	profileData := struct {
		User       models.User                 `json:"user"`
		Posts      []models.Post               `json:"user_posts"`
		LikedPosts []models.PostWithCategories `json:"liked_posts"`
	}{
		User:       userData,
		Posts:      userPosts,
		LikedPosts: userLikedPosts,
	}

	posts.WriteJSON(w, http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "Profile data retrieved",
		Data:    profileData,
	})
}
