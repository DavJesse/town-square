package handlers

// import (
// 	"fmt"
// 	"html/template"
// 	"net/http"

// 	"forum/database"
// 	"forum/models"
// )

// func ForumHandler(w http.ResponseWriter, r *http.Request) {
// 	// Fetch posts from the database
// 	posts, err := database.GetAllPosts()
// 	if err != nil {
// 		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
// 		return
// 	}

// 	// Load the HTML template
// 	tmpl, err := template.ParseFiles("web/templates/index.html")
// 	if err != nil {
// 		http.Error(w, "Failed to load template", http.StatusInternalServerError)
// 		return
// 	}

// 	// Execute the template with the posts data
// 	data := struct {
// 		Posts []models.PostWithUsername
// 	}{
// 		Posts: posts,
// 	}

// 	for i, post := range posts {
// 		fmt.Printf("Uid#%d count is %d\n", i, post.LikesCount)
// 		fmt.Printf("Uid#%d count is %d\n", i, post.DislikesCount)
// 	}

// 	tmpl.Execute(w, data)
// }
