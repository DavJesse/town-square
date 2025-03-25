package posts

import (
	"encoding/json"
	"log"
	"net/http"

	errors "forum/handlers/errors"

	"forum/database"
	"forum/models"
)

func PostDisplay(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	var err error
	session, loggedIn := database.IsLoggedIn(r)

	// Retrieve user data if logged in
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			log.Printf("Error getting user: %v\n", err) // Add error logging
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	postID := r.URL.Query().Get("pid")

	postData, err := database.GetPostByUUID(postID)
	if err != nil {
		errors.NotFoundHandler(w, r)
		log.Println("Error getting post data: ", err)
		return
	}

	// Prepare data for JSON response
	data := struct {
		PostData models.PostWithCategories
		IsLogged bool
		ProfPic  string
	}{
		PostData: postData,
		IsLogged: loggedIn,
		ProfPic:  userData.Image,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
