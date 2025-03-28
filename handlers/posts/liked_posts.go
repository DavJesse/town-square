package posts

import (
	"html/template"
	"log"
	"net/http"

	errors "forum/handlers/errors"
	utils "forum/utils"

	"forum/database"
	"forum/models"
)

type TemplateData struct {
	IsLogged bool
	Posts    []models.PostWithCategories
	ProfPic  string
}

func ShowLikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.MethodNotAllowedHandler(w, r)
		log.Println("METHOD ERROR: method not allowed")
		return
	}

	userID, _, err := database.GetUserData(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve user information
	user, _ := database.GetUserbyID(userID)

	likedPosts, err := database.GetLikedPostsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve liked posts", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		IsLogged: true,
		Posts:    likedPosts,
		ProfPic:  user.Image,
	}

	// Parse template with function to replace '\n' with '<br>'
	tmpl := template.Must(template.New("liked_posts.html").Funcs(template.FuncMap{
		"replaceNewlines": utils.ReplaceNewlines,
	}).ParseFiles("./web/templates/liked_posts.html"))

	// Render the liked posts page
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}
}
