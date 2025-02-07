package comments

import (
	"fmt"
	"net/http"

	"forum/database"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	err = database.LikeComment(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to like comment", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/posts/display?pid=%s", postID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
