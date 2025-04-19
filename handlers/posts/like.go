package posts

import (
	"net/http"

	"forum/database"
	"forum/models"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, models.LikeResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	r.ParseForm()
	postID := r.FormValue("post-id")
	userID, _, err := database.GetUserData(r)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, models.LikeResponse{
			Success: false,
			Message: "Please login first",
		})
		return
	}

	err = database.LikePost(userID, postID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, models.LikeResponse{
			Success: false,
			Message: "Failed to like post",
		})
		return
	}

	// Get updated likes count
	likesCount, err := database.GetPostLikesCount(postID)
	if err != nil {
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
