package comments

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	errors "forum/handlers/errors"
)

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.MethodNotAllowedHandler(w)
		log.Printf("METHOD ERROR: method not allowed")
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
		errors.InternalServerErrorHandler(w)
		fmt.Printf("DISLIKE ERROR: %v", err)
		return
	}

	// Redirect to the posts display page with the postID as a query parameter
	redirectURL := fmt.Sprintf("/posts/display?pid=%s", postID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
