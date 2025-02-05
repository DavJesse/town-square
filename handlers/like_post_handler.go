package handlers

import (
	"fmt"
	"net/http"

	"forum/database"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	postID := r.FormValue("post-id")
	userID := 1 // Replace with actual logged-in user ID

	err := database.LikePost(userID, postID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
