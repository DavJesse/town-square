package handlers

import (
	"fmt"
	"net/http"

	"forum/database"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	// take the contents from the form
	// call the database function to insert a comment
	// redirect to /posts/display?pid={{.UUID}}

	// Only allow POST requests for submitting a comment
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data (assuming the form contains a comment and post UUID)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Retrieve the comment text and post UUID from the form
	commentText := r.FormValue("comment")
	postUUID := r.FormValue("postUUID")
	// userID := r.FormValue("userID") // Assuming userID is passed in the form

	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// fmt.Println("TEXT: ", commentText)

	// Ensure the comment is not empty
	if commentText == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	// // Validate the UUID format (basic validation, adjust according to your needs)
	// if !utils.IsValidUUID(postUUID) {
	// 	http.Error(w, "Invalid post UUID", http.StatusBadRequest)
	// 	return
	// }

	// Assuming you convert the userID from string to int
	// userIDInt, err := strconv.Atoi(userID)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	// Call the CreateComment function to insert the comment into the database
	_, err = database.CreateComment(userID, postUUID, commentText)
	if err != nil {
		http.Error(w, "Error inserting comment into database", http.StatusInternalServerError)
		return
	}

	// Redirect to the post's display page with the post UUID
	http.Redirect(w, r, fmt.Sprintf("/posts/display?pid=%s", postUUID), http.StatusSeeOther)
}
