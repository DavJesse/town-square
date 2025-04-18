// handlers/routes.go
package routes

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/handlers/auth"
	"forum/handlers/comments"
	"forum/handlers/errors"
	"forum/handlers/messages"
	"forum/handlers/middleware"
	"forum/handlers/misc"
	"forum/handlers/posts"
	users "forum/handlers/users"
)

func RegisterRoutes() {
	// Index & static
	http.HandleFunc("/", posts.IndexRouteHandler)
	http.HandleFunc("/api/index-data", posts.GetIndexData)
	http.HandleFunc("/static/", misc.Static)

	// Authentication
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/register", auth.RegistrationHandler)
	http.HandleFunc("/logout", auth.Logout)

	// Posts
	http.HandleFunc("/posts", posts.Posts)
	http.HandleFunc("/posts/display", posts.PostDisplay)
	http.HandleFunc("/categories", posts.GetCategories)
	http.HandleFunc("/categories/", posts.SingleCategoryPosts)
	http.HandleFunc("/search", posts.Search)
	http.HandleFunc("/liked-posts", posts.ShowLikedPosts)
	http.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(posts.PostCreate)))
	http.Handle("/posts/like", middleware.AuthMiddleware(http.HandlerFunc(posts.LikePost)))
	http.Handle("/posts/dislike", middleware.AuthMiddleware(http.HandlerFunc(posts.DislikePost)))

	// Comments
	http.Handle("/comment", middleware.AuthMiddleware(http.HandlerFunc(comments.Comment)))
	http.Handle("/comments/like", middleware.AuthMiddleware(http.HandlerFunc(comments.LikeCommentHandler)))
	http.Handle("/comments/dislike", middleware.AuthMiddleware(http.HandlerFunc(comments.DislikeCommentHandler)))

	// Users
	http.HandleFunc("GET /profile", users.ProfileRouteHandle)
	http.HandleFunc("Get /api/profile-data", users.GetProfileData)
	http.HandleFunc("/api/users", database.GetUsers)
	http.HandleFunc("/api/users/", database.GetUserById)

	// Messages (WebSocket & API)
	messages.NewMessageHub()
	http.Handle("/ws", middleware.AuthMiddleware(http.HandlerFunc(messages.ServeWS)))
	http.HandleFunc("/api/messages/conversations", messages.GetConversationsHandler)
	http.HandleFunc("/api/messages/{userId}", messages.GetMessagesHandler)
	http.HandleFunc("/chat", serveChatPage)

	// Errors
	http.HandleFunc("/error", errors.ErrorHandler)
}

func serveChatPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/chat.html")
}

func WSProcess(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/templates/chat.html"))

	tmpl.Execute(w, nil)
}
