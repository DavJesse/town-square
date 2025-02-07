package posts

import (
	"fmt"
	"net/http"

	"forum/database"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	postID := r.FormValue("post-id")
	// userID := 1 // Replace with actual logged-in user ID

	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = database.LikePost(userID, postID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
