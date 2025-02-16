package auth

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/errors"
	"forum/models"
	"forum/utils"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Construct absolute path to register.html
	tmplPath, err := utils.GetTemplatePath("register.html")
	if err != nil {
		log.Printf("TEMPLATE AVAILABILITY ERROR: %v", err)
		errors.NotFoundHandler(w)
		return
	}

	// Render html template
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("TEMPLATE PARSING ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	// Render registration form when method is GET
	if r.Method == "GET" {
		ExecuteTemplate(w, tmpl)
		return
	}

	// Logout previous session
	err = LogOutSession(w, r)
	if err != nil {
		log.Printf("LOG OUT ERROR: %v", err)
		errors.InternalServerErrorHandler(w)
		return
	}

	// Handle non-GET and non-POST requests
	if r.Method != "POST" {
		log.Println("REQUEST ERROR: bad request")
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
	if err := r.ParseMultipartForm(20 << 20); err != nil { //max size: 2MB
		ParseAlertMessage(w, tmpl, "File upload too large or invalid form data")
		return
	}

	// Validate email format
	if !utils.ValidEmail(EscapeFormSpecialCharacters(r, "email")) {
		ParseAlertMessage(w, tmpl, "Invalid email format")
		return
	}

	// Check if email or username is taken
	existingUser, _ := database.GetUserByEmailOrUsername(EscapeFormSpecialCharacters(r, "email"), EscapeFormSpecialCharacters(r, "username"))
	if existingUser.Username != "" {
		ParseAlertMessage(w, tmpl, fmt.Sprintf("%s taken!", r.FormValue("email")))
		return
	}

	// Extract form data
	user.Email = EscapeFormSpecialCharacters(r, "email")
	user.Username = EscapeFormSpecialCharacters(r, "username")
	password := EscapeFormSpecialCharacters(r, "password")
	confirmPassword := EscapeFormSpecialCharacters(r, "confirm_password")
	user.Bio = EscapeFormSpecialCharacters(r, "bio")

	// Validate passwords
	if password != confirmPassword {
		ParseAlertMessage(w, tmpl, "Passwords do not match")
		return
	}

	// Check password strength
	if err = utils.PasswordStrength(password); err != nil {
		ParseAlertMessage(w, tmpl, err.Error())
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
			ParseAlertMessage(w, tmpl, "Invalid file type. Only PNG and JPG images are allowed")
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
	}

	// Create new user in the database
	_, err = database.CreateNewUser(user)
	if err != nil {
		log.Printf("DATABASE ERROR: %v", err)
		ParseAlertMessage(w, tmpl, "Registration failed, try again")
		return
	}

	// Redirect user to login page
	if w.Header().Get("Content-Type") == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func EscapeFormSpecialCharacters(r *http.Request, elementName string) string {
	return html.EscapeString(r.FormValue(elementName))
}
