package posts

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"

	"forum/database"
	errors "forum/handlers/errors"
	"forum/models"
	"forum/utils"
)

type PostResponse struct {
	Categories []models.Category `json:"categories"`
	IsLogged   bool              `json:"is_logged"`
	ProfPic    string            `json:"prof_pic"`
	Message    string            `json:"message"`
	Code       int               `json:"code,omitempty"`
}

func PostCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetPostCreate(w, r)
	case http.MethodPost:
		handlePostPostCreate(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, PostResponse{
			Message: "METHOD ERROR: method not allowed",
			Code:    http.StatusMethodNotAllowed,
		})
		log.Println("METHOD ERROR: method not allowed")
	}
}

func handleGetPostCreate(w http.ResponseWriter, r *http.Request) {
	session, loggedIn := database.IsLoggedIn(r)

	userData, _ := database.GetUserbySessionID(session.SessionID)

	categories, err := database.FetchCategories()
	if err != nil {
		errors.InternalServerErrorHandler(w, r)
		return
	}

	writeJSON(w, http.StatusOK, PostResponse{
		Categories: categories,
		IsLogged:   loggedIn,
		ProfPic:    userData.Image,
	})
}

func handlePostPostCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(26 * 1024 * 1024); err != nil {
		http.Error(w, "Form size too large", http.StatusRequestEntityTooLarge)
		log.Println("INFO: Client form exceeds 26MB")
		return
	}

	title := r.FormValue("title")
	content := html.EscapeString(r.FormValue("content"))
	categoryIDs := r.Form["categories"]
	if len(categoryIDs) == 0 {
		categoryIDs = []string{"1"}
	}

	var filename string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		if handler.Size > 20*1024*1024 {
			writeJSON(w, http.StatusBadRequest, PostResponse{Message: "image exceeds 20MB"})
			return
		}

		buf := make([]byte, 512)
		if _, err := file.Read(buf); err != nil {
			writeJSON(w, http.StatusInternalServerError, PostResponse{Message: "ERROR: reading file"})
			return
		}
		file.Seek(0, io.SeekStart)

		fileType := http.DetectContentType(buf)
		if !map[string]bool{
			"image/png":  true,
			"image/jpeg": true,
			"image/gif":  true,
		}[fileType] {
			writeJSON(w, http.StatusBadRequest, PostResponse{Message: "Invalid file type. Only GIF, PNG, and JPG images are allowed."})
			return
		}

		filename, err = utils.SaveImage(fileType, file, utils.MEDIA)
		if err != nil {
			log.Println("ERROR: Failed to save image")
			writeJSON(w, http.StatusInternalServerError, PostResponse{Message: "ERROR: Failed to save image"})
			return
		}
	}

	var categoryIDsInt []int
	for _, idStr := range categoryIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("INFO: Invalid Category ID")
			writeJSON(w, http.StatusBadRequest, PostResponse{Message: "INFO: Invalid category ID"})
			return
		}
		categoryIDsInt = append(categoryIDsInt, id)
	}

	if err := database.ValidateCategories(categoryIDsInt); err != nil {
		log.Printf("INFO: Invalid category %v", err)
		writeJSON(w, http.StatusBadRequest, PostResponse{Message: "INFO: Invalid Category"})
		return
	}

	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if _, err := database.CreatePostWithCategories(userID, title, content, filename, categoryIDsInt); err != nil {
		log.Printf("ERROR: Failed to create post: %v", err)
		writeJSON(w, http.StatusInternalServerError, PostResponse{Message: fmt.Sprintf("ERROR: Failed to create post: %v", err)})
		return
	}

	writeJSON(w, http.StatusOK, PostResponse{Message: "Post created successfully"})
}

func writeJSON(w http.ResponseWriter, status int, data PostResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
