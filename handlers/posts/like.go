package posts

import (
	"log"
	"net/http"

	"forum/database"
	errors "forum/handlers/errors"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.MethodNotAllowedHandler(w, r)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	r.ParseForm()
	postID := r.FormValue("post-id")
	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = database.LikePost(userID, postID)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	// Redirect to the previous page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
