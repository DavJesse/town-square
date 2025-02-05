package handlers

import (
	"fmt"
	"forum/database"
	"fmt"
	"forum/database"
	"net/http"
	"html/template"
	 "log"
	 //"forum/models"
	"html/template"
	 "log"
	 //"forum/models"
)

// ViewUserProfile handler
func ViewUserProfile(w http.ResponseWriter, r *http.Request) {
    cookieExists, err := HasCookie(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        fmt.Println("Redirected to login")
        return
    }

    if !cookieExists {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Error(w, "Failed to get cookie", http.StatusUnauthorized)
        return
    }

    userData, err := database.GetUserbySessionID(cookie.Value)
   // fmt.Printf("UserData retrieved: %+v\n", userData)  // Add debug logging
    if err != nil {
        fmt.Printf("Error getting user: %v\n", err)  // Add error logging
        http.Error(w, "Session invalid", http.StatusUnauthorized)
        return
    }

    // // Create a data structure for the template
    // data := struct {
    //     User models.User
    // }{
    //     User: userData,
    // }

	// Render the template with data
	path,err:=GetTemplatePath("profile.html")
	if err!=nil{
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
