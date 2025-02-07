package auth

import (
	"fmt"
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
		errors.InternalServerErrorHandler(w)
		log.Println("Could not find template file: ", err)
		return
	}

	// Render html template
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Unable to parse registration template", http.StatusInternalServerError)
		return
	}

	// Render registration form when method is GET
	if r.Method == "GET" {
		err := tmpl.Execute(w, nil)
		if err != nil {
			errors.InternalServerErrorHandler(w)
			log.Println("Could not render registration template: ", err)
			return
		}
		return
	}

	// Logout previous session
	err = LogOutSession(w, r)
	if err != nil {
		errors.InternalServerErrorHandler(w)
		log.Println("Error logging out session: ", err)
		return
	}

	// Handle non-GET and non-POST requests
	if r.Method != "POST" {
		errors.BadRequestHandler(w)
		log.Println("RegistrationHandler ERROR: Bad request")
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		errors.BadRequestHandler(w)
		log.Println("Invalid form submission", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil { // 20MB max
		ParseAlertMessage(w, tmpl, "File upload too large or invalid form data")
		return
	}

	// Validate email format
	if !utils.ValidEmail(r.FormValue("email")) {
		ParseAlertMessage(w, tmpl, "Invalid email format")
		return
	}

	// Check if email or username is taken
	existingUser, _ := database.GetUserByEmailOrUsername(r.FormValue("email"), r.FormValue("username"))
	if existingUser.Username != "" {
		ParseAlertMessage(w, tmpl, fmt.Sprintf("%s taken!", r.FormValue("email")))
		return
	}

	// Extract form data
	user.Email = r.FormValue("email")
	user.Username = r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	user.Bio = r.FormValue("bio")

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

	// If provided, upload image parsed as profile picture
	file, handler, err := r.FormFile("image")
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

	// save the image to disk
	filename, err := utils.SaveImage(fileType, file, utils.IMAGES)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update the name of the profile image
	user.Image = filename

	// Create new user in the database
	_, err = database.CreateNewUser(user)
	if err != nil {
		ParseAlertMessage(w, tmpl, "registration in failed, try again")
		log.Println("Error creating user")
		return
	}

	// Redirect user to login page
	if w.Header().Get("Content-Type") == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
