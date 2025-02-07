package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/database"
	"forum/handlers"
	postHandlers "forum/handlers/posts"
	"forum/utils"
)

func init() {
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

	// Candle hundler functions
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/static/", handlers.StaticHandler)
	http.HandleFunc("/success", handlers.SuccessHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/forgot-password", handlers.ForgotPasswordHandler)
	http.HandleFunc("/register", handlers.RegistrationHandler)
	// RESTORE // http.Handle("/posts/create", middleware.AuthMiddleware(http.HandlerFunc(postHandlers.PostCreate)))
	http.HandleFunc("/posts/create", postHandlers.PostCreate)
	http.HandleFunc("/posts/display", postHandlers.PostDisplay)
	// User Profile routes
	http.HandleFunc("GET /profile", handlers.ViewUserProfile)
	// http.HandleFunc("GET /user/update", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateUserProfile))) // Protected

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		user, _ := database.GetUserByEmailOrUsername("toni", "toni")
		fmt.Println("User: ", user)
	})

	http.HandleFunc("/posts/like", handlers.LikePostHandler)
	http.HandleFunc("/posts/dislike", handlers.DislikePostHandler)
	http.HandleFunc("/comments/like", handlers.LikeCommentHandler)
	http.HandleFunc("/comments/dislike", handlers.DislikeCommentHandler)
	http.HandleFunc("/comment", handlers.Comment)
	http.HandleFunc("/categories", handlers.CategoriesPageHandler)
	http.HandleFunc("/categories/", handlers.SingeCategoryPosts)

	// Inform user initialization of server
	log.Printf("Server runing on http://localhost%s\n", port)

	// Start the server, handle emerging errors
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Failed to start server: ", err)
		return
	}
}
