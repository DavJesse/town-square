package posts

import (
	"net/http"

	"forum/database"
)

func DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
