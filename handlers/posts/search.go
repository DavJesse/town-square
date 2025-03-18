package posts

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	errors "forum/handlers/errors"
	"forum/models"
)

// SearchHandler handles search requests and returns JSON response
func Search(w http.ResponseWriter, r *http.Request) {
	var user models.User

	session, logged := database.IsLoggedIn(r)

	if logged {
		user, _ = database.GetUserbySessionID(session.SessionID)
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		errors.BadRequestHandler(w, r)
		log.Println("Query parameter is missing")
		return
	}

	posts, err := database.SearchPosts(query)
	if err != nil {
		log.Println("Error searching posts:", err)
		errors.InternalServerErrorHandler(w, r)
		return
	}

	// Create a response structure
	data := struct {
		Posts    []models.Post `json:"posts"`
		IsLogged bool          `json:"is_logged"`
		ProfPic  string        `json:"prof_pic"`
	}{
		Posts:    posts,
		IsLogged: logged,
		ProfPic:  user.Image,
	}

	// Set response header to indicate JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode data into JSON and send it as the response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding JSON:", err)
		errors.InternalServerErrorHandler(w, r)
	}
}
