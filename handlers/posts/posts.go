package posts

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
)

func Posts(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	var err error
	session, loggedIn := database.IsLoggedIn(r)
	// Retrieve user data if logged in
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			log.Printf("Error getting user: %v\n", err) // Add error logging
			// Return JSON response with redirect
			jsonResponse := models.Response{
				Code:     http.StatusSeeOther,
				Message:  "Error getting user",
				Redirect: "/login",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusSeeOther)
			json.NewEncoder(w).Encode(jsonResponse)
			return
		}
	} else {
		// not logged in -> redirect to login page
		errors.RedirectHandler(w, "/login")
		return

	}

	// Fetch posts from the database
	posts, err := database.GetAllPosts()
	if err != nil {
		// Use existing error handler for internal server error
		errors.InternalServerErrorHandler(w, r)
		return
	}
	categories, err := database.FetchCategories()
	if err != nil {
		// Use existing error handler for internal server error
		errors.InternalServerErrorHandler(w, r)
		return
	}

	data := struct {
		Posts      []models.PostWithUsername
		Categories []models.Category
		IsLogged   bool   `json:"isLogged"`
		ProfPic    string `json:"ProfPic"`
	}{
		Posts:      posts,
		Categories: categories,
		IsLogged:   loggedIn,
		ProfPic:    userData.Image,
	}

	jsonResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}
