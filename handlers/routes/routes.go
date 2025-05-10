// handlers/routes.go
package routes

import (
	"net/http"

	"forum/handlers/auth"
	"forum/handlers/comments"
	"forum/handlers/errors"
	messages "forum/handlers/messages"
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
	http.Handle("/posts/like", middleware.AuthMiddleware(http.HandlerFunc(posts.LikePostHandler)))
	http.Handle("/posts/dislike", middleware.AuthMiddleware(http.HandlerFunc(posts.DislikePostHandler)))

	// Comments
	http.Handle("/comment", middleware.AuthMiddleware(http.HandlerFunc(comments.Comment)))
	http.Handle("/comments/like", middleware.AuthMiddleware(http.HandlerFunc(comments.LikeCommentHandler)))
	http.Handle("/comments/dislike", middleware.AuthMiddleware(http.HandlerFunc(comments.DislikeCommentHandler)))

	// Users
	http.HandleFunc("GET /profile", users.ProfileRouteHandle)
	http.HandleFunc("/api/profile-data", users.GetProfileData)

	// messages
	http.Handle("/ws", middleware.AuthMiddleware(http.HandlerFunc(messages.HandleWebSocket)))
	http.Handle("/users", middleware.AuthMiddleware(messages.GetAllUsers()))
	http.Handle("/messages", middleware.AuthMiddleware(messages.GetMessagesHandler()))
	http.Handle("POST /send_message", middleware.AuthMiddleware(messages.SendMessageHTTPHandler()))

	// Errors
	http.HandleFunc("/error", errors.ErrorHandler)
}
