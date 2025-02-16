package posts

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"

	errors "forum/handlers/errors"

	"forum/database"
	"forum/models"
	utils "forum/utils"
)

// Handler for serving the form and handling form submission
func PostCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var loggedIn bool
		session, lIn := database.IsLoggedIn(r)
		if lIn {
			loggedIn = true
		}

		// Retrieve user data
		userData, _ := database.GetUserbySessionID(session.SessionID)

		// fetch categories from the database
		categories, err := database.FetchCategories()
		if err != nil {
			errors.InternalServerErrorHandler(w)
			return
		}

		data := struct {
			Categories []models.Category
			IsLogged   bool
			ProfPic    string
		}{
			Categories: categories,
			IsLogged:   loggedIn,
			ProfPic:    userData.Image,
		}

		tmpl := template.Must(template.ParseFiles("./web/templates/posts_create.html"))
		tmpl.Execute(w, data)

	case http.MethodPost:
		if err := r.ParseMultipartForm(20 << 20); err != nil { // 20MB max
			http.Error(w, "File upload too large or invalid form data", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := html.EscapeString(r.FormValue("content"))
		categoryIDs := r.Form["categories"] // Get selected category IDs

		// if no categories are selected, default to category ID 1
		if len(categoryIDs) == 0 {
			categoryIDs = append(categoryIDs, "1")
		}

		// handle the uploaded file if present
		var filename string
		file, handler, err := r.FormFile("image")
		if err == nil { // Only process the file if it's uploaded
			defer file.Close()

			// Validate the file extension type and size
			allowedTypes := map[string]bool{
				"image/png":  true,
				"image/jpeg": true,
			}
			fileType := handler.Header.Get("Content-Type")
			if !allowedTypes[fileType] {
				http.Error(w, "Invalid file type. Only PNG and JPG images are allowed.", http.StatusBadRequest)
				return
			}

			// Save the image to disk
			filename, err = utils.SaveImage(fileType, file, utils.MEDIA)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Convert category IDs from strings to integers
		var categoryIDsInt []int
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}
			categoryIDsInt = append(categoryIDsInt, id)
		}

		// Validate that the selected categories exist in the database
		if err := database.ValidateCategories(categoryIDsInt); err != nil {
			http.Error(w, fmt.Sprintf("Invalid category: %v", err), http.StatusBadRequest)
			return
		}

		// Get user data
		userID, _, err := database.GetUserData(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Create the post with categories
		_, err = database.CreatePostWithCategories(userID, title, content, filename, categoryIDsInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create post: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return

	default:
		errors.MethodNotAllowedHandler(w)
		log.Println("METHOD ERROR: method not allowed")
	}
}
