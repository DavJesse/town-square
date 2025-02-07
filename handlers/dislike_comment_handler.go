package handlers

import (
	"fmt"
	"net/http"

	"forum/database"
)

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	commentID := r.FormValue("comment-id")
	postID := r.FormValue("post-id")
	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = database.DislikeComment(userID, commentID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to dislike comment", http.StatusInternalServerError)
		return
	}

	// Redirect to the posts display page with the postID as a query parameter
	redirectURL := fmt.Sprintf("/posts/display?pid=%s", postID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
