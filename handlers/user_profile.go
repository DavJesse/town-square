package handlers

import (
	"forum/database"
	"net/http"
	"html/template"
	 "log"
	 "forum/models"
)

// ViewUserProfile handler
func ViewUserProfile(w http.ResponseWriter, r *http.Request) {
    cookieExists, cookie, err := HasCookie(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if !cookieExists {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    userData, err := database.GetUserbySessionID(cookie.Value)
   // fmt.Printf("UserData retrieved: %+v\n", userData)  // Add debug logging
    if err != nil {
        log.Printf("Error getting user: %v\n", err)  // Add error logging
		http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

	UserPosts,err:=database.GetUserPostsbyUserID(userData.ID)
	if err != nil {
        log.Printf("Error getting posts: %v\n", err)  // Add error logging
    }

	// Render the template with data
	path,err:=GetTemplatePath("profile.html")
	if err!=nil{
		log.Println("Error getting template path")
	}

  // Combine user data and user posts into a single struct
  profileData := struct{
	User models.User
	Posts []models.Post
  }{
	User:  userData,
	Posts: UserPosts,
}

    tmpl, err := template.ParseFiles(path)
    if err != nil {
        http.Error(w, "Failed to load profile template", http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, profileData); err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

// func UpdateUserProfile(){
// 	// Update user profile
// }

