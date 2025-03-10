package auth

import (
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
	"forum/utils"
)

const MaxUploadSize = 20 * 1024 * 1024

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Logout previous session
	err := LogOutSession(w, r)
	if err != nil {
		log.Printf("LOG OUT ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	// Serve page at get request
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./web/templates/index.html")
		return
	}

	// Handle non-GET and non-POST requests
	if r.Method != http.MethodPost {
		log.Println("REQUEST ERROR: bad request")
		errors.BadRequestHandler(w)
		return
	}

	// Check content type
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		errors.BadRequestHandler(w)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		log.Println("REQUEST ERROR: bad request")
		errors.BadRequestHandler(w)
		return
	}

	// Validate image uploaded from form
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil { // max size: 20MB
		ParseAlertMessage(w, "File upload too large or invalid form data")
		return
	}

	// Validate email format
	if !utils.ValidEmail(EscapeFormSpecialCharacters(r, "email")) {
		ParseAlertMessage(w, "Invalid email format")
		return
	}

	// Check if email is taken
	existingUser, _ := database.GetUserByEmailOrUsername(EscapeFormSpecialCharacters(r, "email"), EscapeFormSpecialCharacters(r, "username"))
	if existingUser.Username != "" {
		ParseAlertMessage(w, "the email/username taken")
		return
	}

	// Extract form data
	user.FirstName = EscapeFormSpecialCharacters(r, "first_name")
	user.LastName = EscapeFormSpecialCharacters(r, "last_name")
	user.Username = EscapeFormSpecialCharacters(r, "username")
	user.Age, _ = strconv.Atoi(EscapeFormSpecialCharacters(r, "age"))
	user.Gender = EscapeFormSpecialCharacters(r, "gender")
	user.Email = EscapeFormSpecialCharacters(r, "email")
	password := EscapeFormSpecialCharacters(r, "password")
	confirmPassword := EscapeFormSpecialCharacters(r, "confirm_password")
	user.Bio = EscapeFormSpecialCharacters(r, "bio")

	// Validate passwords
	if password != confirmPassword {
		ParseAlertMessage(w, "Passwords do not match")
		return
	}

	// Check password strength
	if err = utils.PasswordStrength(password); err != nil {
		ParseAlertMessage(w, err.Error())
		return
	}

	user.Password = password  // set password
	utils.Passwordhash(&user) // Hash password

	// Check if a file is uploaded for the profile image
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
			ParseAlertMessage(w, "Invalid file type. Only PNG and JPG images are allowed")
			return
		}

		// Save the image to disk
		filename, err = utils.SaveImage(fileType, file, utils.IMAGES)
		if err != nil {
			log.Printf("IMAGE SAVING ERROR: %v", err)
			errors.InternalServerErrorHandler(w)
			return
		}

		// Update the name of the profile image
		user.Image = filename

		// Capture errors due to missing files
	} else if err == http.ErrMissingFile {
		log.Printf("FILE UPLOAD ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	// Create new user in the database
	_, err = database.CreateNewUser(user)
	if err != nil {
		log.Printf("DATABASE ERROR: %v", err)
		ParseAlertMessage(w, "Registration failed, try again")
		return
	}

	// Redirect user to login page
	if w.Header().Get("Content-Type") == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func EscapeFormSpecialCharacters(r *http.Request, elementName string) string {
	return html.EscapeString(r.Form.Get(elementName))
}
