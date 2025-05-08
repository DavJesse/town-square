package posts

import (
	"log"
	"net/http"

	"forum/database"
)

var postID string

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Block non-POST and non-GET requests
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	if r.Method == "POST" {
		// Capture post ID
		r.ParseForm()
		postID := r.FormValue("post-id")

		// Capture  user data while blocking unauthorized users
		userID, _, err := database.GetUserData(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = database.LikePost(userID, postID)
		if err != nil {
			log.Printf("DATABASE ERROR: Failed to Log Like in Database, %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Send updated like and dislike counts at GET request
	if r.Method == "GET" {
		postEngagement, err := database.GetPostEngagementCount(postID)
		if err != nil {
			log.Println("DATABASE ERROR: Failed to Retrieve Likes Count")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		postID = "" // Reset postID
		WriteJSON(w, http.StatusOK, postEngagement)
	}
}
