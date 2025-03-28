package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/models"
)

// ViewUserProfile handler
func ViewUserProfile(w http.ResponseWriter, r *http.Request) {
	session, logged := database.IsLoggedIn(r)
	if !logged {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userData, err := database.GetUserbySessionID(session.SessionID)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Message: "Profile data retrieved",
		Data:    profileData,
	})
}
