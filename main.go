package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/database"

	auth "forum/handlers/auth"
	comments "forum/handlers/comments"
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
}

func main() {
	portStr := utils.Port() // get the port to use to start the server
	port := fmt.Sprintf(":%d", portStr)

	// will postpone the closure of the database handler created by init/0 function to when main/0 exits
	defer database.Close()

	// Restrict arguments parsed
	if len(os.Args) != 1 {
		log.Println("Too many arguments")
		log.Println("Usage: go run main.go")
		return
	}

	// authentication
	http.HandleFunc("/", posts.Index)
	http.HandleFunc("/static/", misc.Static)
	http.HandleFunc("/login", auth.Login)
	// http.HandleFunc("/forgot-password", auth.ForgotPassword) // Unmute when retrieval logic is implemented
	http.HandleFunc("/register", auth.Registration)
	http.HandleFunc("/logout", auth.Logout)

	// users
	http.HandleFunc("GET /profile", users.ViewUserProfile)
	// http.HandleFunc("GET /user/update", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateUserProfile))) // Protected

	// posts
	http.HandleFunc("/posts/display", posts.PostDisplay)
	http.HandleFunc("/categories", posts.CategoriesPage)
	http.HandleFunc("/categories/", posts.SingeCategoryPosts)
	http.HandleFunc("/search", posts.Search)
	http.HandleFunc("/liked-posts", posts.ShowLikedPosts)

	http.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(posts.PostCreate)))
	http.Handle("/posts/like", middleware.AuthMiddleware(http.HandlerFunc(posts.LikePost)))
	http.Handle("/posts/dislike", middleware.AuthMiddleware(http.HandlerFunc(posts.DislikePost)))

	// comments
	http.Handle("/comments/like", middleware.AuthMiddleware(http.HandlerFunc(comments.LikeCommentHandler)))
	http.Handle("/comments/dislike", middleware.AuthMiddleware(http.HandlerFunc(comments.DislikeCommentHandler)))
	http.Handle("/comment", middleware.AuthMiddleware(http.HandlerFunc(comments.Comment)))

	// start the server, handle emerging errors
	fmt.Printf("Server runing on http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Failed to start server: ", err)
		return
	}
}
