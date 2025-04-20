package posts

import (
	"log"
	"net/http"

	"forum/database"
	"forum/models"
)

func DislikePost(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, isLogged := database.IsLoggedIn(r)
	if !isLogged {
		WriteJSON(w, http.StatusUnauthorized, models.DislikeResponse{
			Success: false,
			Message: "Please login first",
		})
		return
	}

	r.ParseForm()
	postID := r.FormValue("post-id")
	userID, _, err := database.GetUserData(r)

	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, models.DislikeResponse{
			Success: false,
			Message: "Please login first",
		})
		return
	}

	err = database.DislikePost(userID, postID)
	if err != nil {
		log.Println("DATABASE ERROR: Failed to Log Dislike in Database")
		WriteJSON(w, http.StatusInternalServerError, models.DislikeResponse{
			Success: false,
			Message: "Failed to dislike post",
		})
		return
	}

	// Get updated likes count
	likesCount, err := database.GetPostDislikesCount(postID)
	if err != nil {
		log.Println("DATABASE ERROR: Failed to Retrieve Dislikes Count")
		WriteJSON(w, http.StatusInternalServerError, models.DislikeResponse{
			Success: false,
			Message: "Failed to get updated dislikes count",
		})
		return
	}

	WriteJSON(w, http.StatusOK, models.DislikeResponse{
		Success:       true,
		DislikesCount: likesCount,
		Message:       "Post liked successfully",
	})
}
