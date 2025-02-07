package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"forum/database"
	"forum/models"
	"forum/utils"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Construct absolute path to register.html
	tmplPath, err := GetTemplatePath("register.html")
	if err != nil {
		InternalServerErrorHandler(w)
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
			InternalServerErrorHandler(w)
			log.Println("Could not render registration template: ", err)
			return
		}
		return
	}

	// Logout previous session
	err = LogOutSession(w, r)
	if err != nil {
		InternalServerErrorHandler(w)
		log.Println("Error logging out session: ", err)
		return
	}

	// Handle non-GET and non-POST requests
	if r.Method != "POST" {
		BadRequestHandler(w)
		log.Println("RegistrationHandler ERROR: Bad request")
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		BadRequestHandler(w)
		log.Println("Invalid form submission", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil { // 20MB max
		ParseAlertMessage(w, tmpl, "File upload too large or invalid form data")
		return
	}

	// Validate email format
	if !ValidEmail(r.FormValue("email")) {
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
	if r.FormValue("image") != "" {
		UploadImage(w, r, &user, tmpl)
	}

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

// Validate email format

func ValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}

func UploadImage(w http.ResponseWriter, r *http.Request, user *models.User, tmpl *template.Template) {
	// Handle the uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		ParseAlertMessage(w, tmpl, "image upload failed")
		log.Panicln("Failed to retrieve the file")
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
		ParseAlertMessage(w, tmpl, "Invalid file type. Only PNG and JPG images are allowed.")
		log.Println("Invalid file type. Only PNG and JPG images are allowed.")
		return
	}

	// Generate a random filename
	randomFileName, err := utils.GenerateRandomName()
	if err != nil {
		ParseAlertMessage(w, tmpl, "image upload failed")
		log.Println("Failed to generate a unique filename")
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
	user.Image = filename

	// Save the file to the media folder
	mediaFolder := "web/static/images"
	if err := os.MkdirAll(mediaFolder, os.ModePerm); err != nil {
		ParseAlertMessage(w, tmpl, "image upload failed")
		log.Println("Failed to create media folder")
		return
	}

	filePath := filepath.Join(mediaFolder, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		ParseAlertMessage(w, tmpl, "image upload failed")
		log.Println("Failed to save the file")
		return
	}

	defer outFile.Close()
	// Copy the file content to the new file
	if _, err := io.Copy(outFile, file); err != nil {
		ParseAlertMessage(w, tmpl, "image upload failed")
		log.Println("Failed to save the file")
		return
	}
}
