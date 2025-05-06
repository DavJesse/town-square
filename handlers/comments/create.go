package comments

import (
	"encoding/json"
	"html"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/posts"
	"forum/models"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %v", r.Method)
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
		log.Printf("Content-type: %v", r.Header.Get("Content-Type"))
		var requestData struct {
			Comment string `json:"comment"`
			PostID  string `json:"postUUID"`
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestData)
		if err != nil {
			posts.WriteJSON(w, http.StatusBadRequest, models.PostResponse{
				Message: "REQUEST ERROR: bad request",
				Code:    http.StatusBadRequest,
			})
			log.Printf("JSON DECODE ERROR: %v", err)
			return
		}

		// Retrieve the comment text and post UUID from the form
		commentText := html.EscapeString(requestData.Comment)
		postID = requestData.PostID

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
