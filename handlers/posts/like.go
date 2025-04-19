package posts

import (
	"log"
	"net/http"

	"forum/database"
	"forum/models"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, isLogged := database.IsLoggedIn(r)
	if !isLogged {
		WriteJSON(w, http.StatusUnauthorized, models.LikeResponse{
			Success: false,
			Message: "Please login first",
		})
		return
	}

	r.ParseForm()
	postID := r.FormValue("post-id")
	userID, _, err := database.GetUserData(r)
	log.Printf("User Id: %v -- Post Id: %v", userID, postID)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, models.LikeResponse{
			Success: false,
			Message: "Please login first",
		})
		return
	}

	err = database.LikePost(userID, postID)
	if err != nil {
		log.Println("DATABASE ERROR: Failed to Log Like in Database")
		WriteJSON(w, http.StatusInternalServerError, models.LikeResponse{
			Success: false,
			Message: "Failed to like post",
		})
		return
	}

	// Get updated likes count
	likesCount, err := database.GetPostLikesCount(postID)
	if err != nil {
		log.Println("DATABASE ERROR: Failed to Retrieve Likes Count")
		WriteJSON(w, http.StatusInternalServerError, models.LikeResponse{
			Success: false,
			Message: "Failed to get updated likes count",
		})
		return
	}

	WriteJSON(w, http.StatusOK, models.LikeResponse{
		Success:    true,
		LikesCount: likesCount,
		Message:    "Post liked successfully",
	})
}
