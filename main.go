package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"forum/database"

	auth "forum/handlers/auth"
	comments "forum/handlers/comments"
	errors "forum/handlers/errors"
	messages "forum/handlers/messages"
	"forum/handlers/middleware"
	"forum/handlers/misc"
	posts "forum/handlers/posts"
	users "forum/handlers/users"
	utils "forum/utils"
)

func init() {
	utils.CreatImagesFolder()
	utils.CreatMediaFolder()
	utils.CreatStorageFolder()

	err := database.Init("storage/forum.db")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// database.StoreMessage(&models.Message{
	// 	Content:    "Hello there, I am raja",
	// 	SenderID:   1,
	// 	ReceiverID: 2,
	// 	Timestamp:  time.Now().UTC(),
	// })
	// database.StoreMessage(&models.Message{
	// 	Content:    "Hi, Mi Markie",
	// 	SenderID:   2,
	// 	ReceiverID: 1,
	// 	Timestamp:  time.Now().UTC(),
	// })
	// database.StoreMessage(&models.Message{
	// 	Content:    "Sup Markie, whats poppin?",
	// 	SenderID:   1,
	// 	ReceiverID: 2,
	// 	Timestamp:  time.Now().UTC(),
	// })
	// database.StoreMessage(&models.Message{
	// 	Content:    "Markie mi a beat them bad",
	// 	SenderID:   2,
	// 	ReceiverID: 1,
	// 	Timestamp:  time.Now().UTC(),
	// })
}

func main() {
	portStr := utils.Port() // get the port to use to start the server
	port := fmt.Sprintf(":%d", portStr)

	// will postpone the closure of the database handler created by init/0 function to when main/0 exits
	defer database.Close()

	// Restrict arguments parsed
	if len(os.Args) != 1 {
		log.Println("Too many arguments")
		log.Println("Usage: go run .")
		return
	}

	// authentication
	http.HandleFunc("/", posts.Index)

	// Create a new WebSocket server instance
	// messages.NewWebSocketServer()
	// http.Handle("/chat", middleware.AuthMiddleware(http.HandlerFunc(messages.WebSocketHandler)))
	// http.HandleFunc("/ws", WSProcess)

	// ============================================================================================================================

	// Create a single MessageHub instance
	messages.NewMessageHub()

	// Message routes
	http.HandleFunc("/api/messages/conversations", messages.GetConversationsHandler)

	http.HandleFunc("/api/messages/{userId}", messages.GetMessagesHandler)

	// WebSocket route
	http.Handle("/ws", middleware.AuthMiddleware(http.HandlerFunc(messages.ServeWS)))
	// User list route (returns registered users)
	http.HandleFunc("/api/users", database.GetUsers)

	http.HandleFunc("/api/users/", database.GetUserById)

	http.HandleFunc("/chat", WSProcess)

	// ======================================================================================================

	http.HandleFunc("/static/", misc.Static)
	http.HandleFunc("/login", auth.LoginHandler)
	// http.HandleFunc("/forgot-password", auth.ForgotPassword) // Unmute when retrieval logic is implemented
	http.HandleFunc("/register", auth.RegistrationHandler)
	http.HandleFunc("/logout", auth.Logout)

	// users
	http.HandleFunc("GET /profile", users.ViewUserProfile)
	// http.HandleFunc("GET /user/update", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateUserProfile))) // Protected

	// posts
	http.HandleFunc("/posts", posts.Posts)
	http.HandleFunc("/posts/display", posts.PostDisplay)
	http.HandleFunc("/categories", posts.GetCategories)
	http.HandleFunc("/categories/", posts.SingleCategoryPosts)
	http.HandleFunc("/search", posts.Search)
	http.HandleFunc("/liked-posts", posts.ShowLikedPosts)

	http.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(posts.PostCreate)))
	http.Handle("/posts/like", middleware.AuthMiddleware(http.HandlerFunc(posts.LikePost)))
	http.Handle("/posts/dislike", middleware.AuthMiddleware(http.HandlerFunc(posts.DislikePost)))

	// comments
	http.Handle("/comments/like", middleware.AuthMiddleware(http.HandlerFunc(comments.LikeCommentHandler)))
	http.Handle("/comments/dislike", middleware.AuthMiddleware(http.HandlerFunc(comments.DislikeCommentHandler)))
	http.Handle("/comment", middleware.AuthMiddleware(http.HandlerFunc(comments.Comment)))

	// errors
	http.HandleFunc("/error", errors.ErrorHandler)

	// start the server, handle emerging errors
	fmt.Printf("Server runing on http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Failed to start server: ", err)
		return
	}
}

func WSProcess(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/templates/chat.html"))

	tmpl.Execute(w, nil)
}
