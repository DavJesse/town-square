package comments

import (
	"fmt"
	"log"
	"net/http"

	errors "forum/handlers/errors"

	"forum/database"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	err = database.LikeComment(userID, commentID)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		fmt.Printf("DISLIKE ERROR: %v", err)
		return
	}

	redirectURL := fmt.Sprintf("/posts/display?pid=%s", postID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
