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

	userPosts, err := database.PostsFilterByUser(userData.ID)
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		// Continue execution but return empty posts if an error occurs
		userPosts = []models.Post{}
	}

	// Combine user data and user posts
	profileData := struct {
		User  models.User   `json:"user"`
		Posts []models.Post `json:"posts"`
	}{
		User:  userData,
		Posts: userPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Message: "Profile data retrieved",
		Data:    profileData,
	})
}
