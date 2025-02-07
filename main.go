package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/database"
	"forum/handlers"
	auth "forum/handlers/auth"
	comments "forum/handlers/comments"
	posts "forum/handlers/posts"
	users "forum/handlers/users"
	"forum/utils"
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
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/static/", handlers.Static)
	http.HandleFunc("/success", auth.Success)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/forgot-password", auth.ForgotPassword)
	http.HandleFunc("/register", auth.Registration)
	http.HandleFunc("/logout", auth.Logout)

	// users
	http.HandleFunc("GET /profile", users.ViewProfile)
	// http.HandleFunc("GET /user/update", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateUserProfile))) // Protected

	// posts
	http.HandleFunc("/posts/create", posts.PostCreate)
	http.HandleFunc("/posts/display", posts.PostDisplay)
	http.HandleFunc("/posts/like", posts.LikePost)
	http.HandleFunc("/posts/dislike", posts.DislikePost)
	http.HandleFunc("/categories", posts.CategoriesPage)
	http.HandleFunc("/categories/", posts.SingeCategoryPosts)
	http.HandleFunc("/search", posts.Search)
	// RESTORE // http.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(postHandlers.PostCreate)))

	// comments
	http.HandleFunc("/comments/like", comments.LikeCommentHandler)
	http.HandleFunc("/comments/dislike", comments.DislikeCommentHandler)
	http.HandleFunc("/comment", comments.Comment)

	// start the server, handle emerging errors
	fmt.Printf("Server runing on http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Failed to start server: ", err)
		return
	}
}
