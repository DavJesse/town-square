package posts

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"forum/database"
	"forum/utils"
)

// Handler for serving the form and handling form submission
func PostCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Fetch categories from the database
		categories, err := database.FetchCategories()
		if err != nil {
			http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
			return
		}

		// Serve the form for creating a post
		tmpl := template.Must(template.ParseFiles("./web/templates/posts_create.html"))
		tmpl.Execute(w, categories)

	case http.MethodPost:
		if err := r.ParseMultipartForm(20 << 20); err != nil { // 20MB max
			http.Error(w, "File upload too large or invalid form data", http.StatusBadRequest)
			return
		}

		// Extract the form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryIDs := r.Form["categories"] // Get selected category IDs

		// Handle the uploaded file
		file, handler, err := r.FormFile("media")
		if err != nil {
			http.Error(w, "Failed to retrieve the file", http.StatusBadRequest)
			return
		}
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

		// Generate a random filename
		randomFileName, err := utils.GenerateRandomName()
		if err != nil {
			http.Error(w, "Failed to generate a unique filename", http.StatusInternalServerError)
			return
		}

		// Determine the file extension based on the MIME type
		var ext string
		switch fileType {
		case "image/png":
			ext = ".png"
		case "image/jpeg":
			ext = ".jpg"
		}

		// Construct the full filename
		filename := randomFileName + ext

		// Save the file to the media folder
		mediaFolder := "web/static/media"
		if err := os.MkdirAll(mediaFolder, os.ModePerm); err != nil {
			http.Error(w, "Failed to create media folder", http.StatusInternalServerError)
			return
		}

		filePath := filepath.Join(mediaFolder, filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save the file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy the file content to the new file
		if _, err := io.Copy(outFile, file); err != nil {
			http.Error(w, "Failed to save the file", http.StatusInternalServerError)
			return
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
