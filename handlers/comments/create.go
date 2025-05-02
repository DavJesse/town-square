package comments

import (
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/auth"
	"forum/handlers/posts"
	"forum/models"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests for submitting a comment
	if !(r.Method == http.MethodPost || r.Method == http.MethodGet) {
		posts.WriteJSON(w, http.StatusMethodNotAllowed, models.PostResponse{
			Message: "METHOD ERROR: method not allowed",
			Code:    http.StatusMethodNotAllowed,
		})
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	var postID string

	if r.Method == http.MethodPost {
		// Parse form data (assuming the form contains a comment and post UUID)
		err := r.ParseForm()
		if err != nil {
			posts.WriteJSON(w, http.StatusBadRequest, models.PostResponse{
				Message: "REQUEST ERROR: bad request",
				Code:    http.StatusBadRequest,
			})
			log.Printf("REQUEST ERROR: %v", err)
			return
		}

		// Retrieve the comment text and post UUID from the form
		commentText := auth.EscapeFormSpecialCharacters(r, "comment")
		postID = r.FormValue("postUUID")

		userID, _, err := database.GetUserData(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Call the CreateComment function to insert the comment into the database
		_, err = database.CreateComment(userID, postID, commentText)
		if err != nil {
			posts.WriteJSON(w, http.StatusInternalServerError, models.PostResponse{
				Message: "Internal Server Error",
				Code:    http.StatusInternalServerError,
			})
			log.Printf("DATABASE ERROR: %v", err)
			return
		}

	}

	// Block access to endpoint if postID is empty
	if postID == "" {
		posts.WriteJSON(w, http.StatusMethodNotAllowed, models.PostResponse{
			Message: "METHOD ERROR: method not allowed",
			Code:    http.StatusMethodNotAllowed,
		})
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	// Fetch updated comments for the post
	comments, err := database.GetCommentsForPost(postID)
	if err != nil {
		posts.WriteJSON(w, http.StatusInternalServerError, models.PostResponse{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		})
		log.Printf("COMMENTS DATABASE ERROR: %v", err)
		return
	}

	// Return the updated comments as JSON
	posts.WriteJSON(w, http.StatusOK, comments)
}
