package comments

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/auth"
	errors "forum/handlers/errors"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	// take the contents from the form
	// call the database function to insert a comment
	// redirect to /posts/display?pid={{.UUID}}

	// Only allow POST requests for submitting a comment
	if r.Method != http.MethodPost {
		errors.MethodNotAllowedHandler(w)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	// Parse form data (assuming the form contains a comment and post UUID)
	err := r.ParseForm()
	if err != nil {
		errors.BadRequestHandler(w)
		log.Printf("REQUEST ERROR: %v", err)
		return
	}

	// Retrieve the comment text and post UUID from the form
	commentText := auth.EscapeFormSpecialCharacters(r, "comment")
	postUUID := r.FormValue("postUUID")

	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Call the CreateComment function to insert the comment into the database
	_, err = database.CreateComment(userID, postUUID, commentText)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("DATABASE ERROR: %v", err)
		return
	}

	// Redirect to the post's display page with the post UUID
	http.Redirect(w, r, fmt.Sprintf("/posts/display?pid=%s", postUUID), http.StatusSeeOther)
}
