package posts

import (
	"html/template"
	"log"
	"net/http"

	utils "forum/utils"

	errors "forum/handlers/errors"

	"forum/database"
	"forum/models"
)

func PostDisplay(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	var err error
	session, loggedIn := database.IsLoggedIn(r)

	_, err = utils.GetTemplatePath("post_display.html")
	if err != nil {
		log.Printf("TEMPLATE AVAILABILITY ERROR: %v", err)
		errors.NotFoundHandler(w)
		return
	}

	// Retrieve user data if logged in
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			log.Printf("Error getting user: %v\n", err) // Add error logging
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// Parse template with function to replace '\n' with '<br>'
	tmpl := template.Must(template.New("post_display.html").Funcs(template.FuncMap{
		"replaceNewlines": utils.ReplaceNewlines,
	}).ParseFiles("./web/templates/post_display.html"))

	postID := r.URL.Query().Get("pid")

	postData, err := database.GetPostByUUID(postID)
	if err != nil {
		log.Println("Error getting post data: ", err)
		return
	}

	// Infuse data to be executed with inquiry if user is logged in
	data := struct {
		PostData models.PostWithCategories
		IsLogged bool
		ProfPic  string
	}{
		PostData: postData,
		IsLogged: loggedIn,
		ProfPic:  userData.Image,
	}

	// fmt.Println("POST: ", PostData)

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		errors.InternalServerErrorHandler(w)
	}
}
