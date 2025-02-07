package posts

import (
	"html/template"
	"log"
	"net/http"

	"forum/database"
)

func PostDisplay(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/post_display.html")
	if err != nil {
		http.Error(w, "Failed to load post template", http.StatusInternalServerError)
		return
	}
	postID := r.URL.Query().Get("pid")
	// fmt.Println("SINGLE PID: ", postID)

	PostData, err := database.GetPostByUUID(postID)
	if err != nil {
		log.Println("Error getting post data: ", err)
		return
	}

	// fmt.Println("POST: ", PostData)

	if err := tmpl.Execute(w, PostData); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
