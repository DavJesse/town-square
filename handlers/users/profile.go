package auth

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers/auth"
	"forum/utils"
)

// ViewUserProfile handler
func ViewProfile(w http.ResponseWriter, r *http.Request) {
	cookieExists, cookie, err := auth.HasCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Redirected to login")
		return
	}

	if !cookieExists {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userData, err := database.GetUserbySessionID(cookie.Value)
	// fmt.Printf("UserData retrieved: %+v\n", userData)  // Add debug logging
	if err != nil {
		log.Printf("Error getting user: %v\n", err) // Add error logging
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Redirected to login")
		return
	}
	// Render the template with data
	path, err := utils.GetTemplatePath("profile.html")
	if err != nil {
		fmt.Println("Error getting template path")
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Failed to load profile template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, userData); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// func UpdateUserProfile(){
// 	// Update user profile
// }
