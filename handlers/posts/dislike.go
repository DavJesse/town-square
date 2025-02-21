package posts

import (
	"log"
	"net/http"

	"forum/database"
	errors "forum/handlers/errors"
)

func DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.MethodNotAllowedHandler(w)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post-id")

	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = database.DislikePost(userID, postID)
	if err != nil {
		http.Error(w, "Failed to dislike post/comment", http.StatusInternalServerError)
		return
	}

	// Redirect to the previous page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
