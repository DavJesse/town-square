package posts

import (
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
	// create template path
	temPath, err := utils.GetTemplatePath("posts_create.html")
	if err != nil {
		errors.NotFoundHandler(w)
	}

	tmpl := template.Must(template.ParseFiles(temPath))

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

	// Capture data to be parsed in html template
	data := struct {
		Categories []models.Category
		IsLogged   bool
		ProfPic    string
		Message    string
	}{
		Categories: categories,
		IsLogged:   loggedIn,
		ProfPic:    userData.Image,
	}

	// Handle allowed methods
	switch r.Method {
	case http.MethodGet: // Parse form
		tmpl.Execute(w, data)

		// Qualify data captured in html form
	case http.MethodPost:
		// Parse the form with a max size (for the entire form, including file and other form data)
		if err := r.ParseMultipartForm(26 * 1024 * 1024); err != nil { // Allow 10MB form data for the image + other fields
			log.Println("INFO: Client form exceeds 26MB")
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

			// Check the file size (if greater than 4MB, return an error)
			if handler.Size > 20*1024*1024 { // 20MB in bytes
				data.Message = "image exceeds 20MB"
				ParseFormMessage(w, tmpl, data)
				return
			}

			// Validate the file extension type
			allowedTypes := map[string]bool{
				"image/png":  true,
				"image/jpeg": true,
				"image/gif":  true,
			}
			fileType := handler.Header.Get("Content-Type")
			if !allowedTypes[fileType] {
				data.Message = "Invalid file type. Only GIF, PNG, and JPG images are allowed."
				ParseFormMessage(w, tmpl, data)
				return
			}

			// Save the image to disk
			filename, err = utils.SaveImage(fileType, file, utils.MEDIA)
			if err != nil {
				log.Println("ERROR: Failed to save image")
				errors.InternalServerErrorHandler(w)
				return
			}
		}

		// Convert category IDs from strings to integers
		var categoryIDsInt []int
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("INFO: Invalid category ID")
				errors.BadRequestHandler(w)
				return
			}
			categoryIDsInt = append(categoryIDsInt, id)
		}

		// Validate that the selected categories exist in the database
		if err := database.ValidateCategories(categoryIDsInt); err != nil {
			log.Printf("INFO: Invalid category %v", err)
			errors.BadRequestHandler(w)
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
			log.Printf("ERROR: Failed to create post: %v", err)
			errors.InternalServerErrorHandler(w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return

	default:
		errors.MethodNotAllowedHandler(w)
		log.Println("METHOD ERROR: method not allowed")
	}
}

func ParseFormMessage(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	// Execute the page
	err := tmpl.Execute(w, data)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Printf("TEMPLATE EXECUTION ERROR: %v", err)
		return
	}
}
