package comments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/posts"
)

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Printf("METHOD ERROR: method not allowed")
		return
	}

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&commentLikeData)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("JSON DECODE ERROR: %v", err)
			return
		}

		commentID := commentLikeData.CommentID

		userID, _, err := database.GetUserData(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = database.DislikeComment(userID, commentID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Printf("DISLIKE ERROR: %v", err)
			return
		}
	}

	// Fetch like and dislike counts at GET request
	if r.Method == http.MethodGet {
		engagement, err := database.GetCommentReactionsCounts(commentLikeData.CommentID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("DATABASE ERROR: %v", err)
			return
		}

		posts.WriteJSON(w, http.StatusOK, engagement)
	}
}
