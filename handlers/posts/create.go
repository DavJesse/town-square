package posts

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"forum/database"
	errors "forum/handlers/errors"
	"forum/models"
	utils "forum/utils"
)

// PostCreate handles post creation requests and returns JSON
func PostCreate(w http.ResponseWriter, r *http.Request) {
	var loggedIn bool
	code := http.StatusOK

	session, lIn := database.IsLoggedIn(r)
	if lIn {
		loggedIn = true
	}

	// Retrieve user data
	userData, _ := database.GetUserbySessionID(session.SessionID)

	// Fetch categories from the database
	categories, err := database.FetchCategories()
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	// Create a response structure
	response := struct {
		Categories []models.Category `json:"categories"`
		IsLogged   bool              `json:"is_logged"`
		ProfPic    string            `json:"prof_pic"`
		Message    string            `json:"message"`
		Code       int               `json:"code,omitempty"`
	}{
		Categories: categories,
		IsLogged:   loggedIn,
		ProfPic:    userData.Image,
		Message:    "",
		Code:       code,
	}

	// Handle allowed methods
	switch r.Method {
	case http.MethodGet:
		// Return the categories and login status as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error encoding JSON:", err)
			errors.InternalServerErrorHandler(w, r)
		}

	case http.MethodPost:
		// Parse the form data
		if err := r.ParseMultipartForm(26 * 1024 * 1024); err != nil { // Allow 26MB form data for the image + other fields
			log.Println("INFO: Client form exceeds 26MB")
			http.Error(w, "Form size too large", http.StatusRequestEntityTooLarge)
			return
		}

		// Capture form values
		title := r.FormValue("title")
		content := html.EscapeString(r.FormValue("content"))
		categoryIDs := r.Form["categories"]

		// Default category if none are selected
		if len(categoryIDs) == 0 {
			categoryIDs = append(categoryIDs, "1")
		}

		// Handle file upload if present
		var filename string
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			// Validate file size
			if handler.Size > 20*1024*1024 {
				response.Message = "image exceeds 20MB"
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			// Validate file type
			allowedTypes := map[string]bool{
				"image/png":  true,
				"image/jpeg": true,
				"image/gif":  true,
			}
			fileType := handler.Header.Get("Content-Type")
			if !allowedTypes[fileType] {
				response.Message = "Invalid file type. Only GIF, PNG, and JPG images are allowed."
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			// Save the image
			filename, err = utils.SaveImage(fileType, file, utils.MEDIA)
			if err != nil {
				response.Message = "ERROR: Failed to save image"
				response.Code = http.StatusInternalServerError
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				log.Println("ERROR: Failed to save image")
				return
			}
		}

		// Convert category IDs to integers
		var categoryIDsInt []int
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				response.Message = "INFO: Invalid category ID"
				response.Code = http.StatusBadRequest
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				log.Println("INFO: Invalid Category ID")
				return
			}
			categoryIDsInt = append(categoryIDsInt, id)
		}

		// Validate categories
		if err := database.ValidateCategories(categoryIDsInt); err != nil {
			response.Message = "INFO: Invalid Category"
			response.Code = http.StatusBadRequest

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			log.Printf("INFO: Invalid category %v", err)
			return
		}

		// Get user data
		userID, _, err := database.GetUserData(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Create post
		_, err = database.CreatePostWithCategories(userID, title, content, filename, categoryIDsInt)
		if err != nil {
			response.Message = fmt.Sprintf("ERROR: Failed to create post: %v", err)
			response.Code = http.StatusInternalServerError

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			log.Printf("ERROR: Failed to create post: %v", err)
			return
		}

		// Send success response
		response.Message = "Post created successfully"
		response.Code = http.StatusOK
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		response.Message = "METHOD ERROR: method not allowed"
		response.Code = http.StatusMethodNotAllowed
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Println("METHOD ERROR: method not allowed")
	}
}
