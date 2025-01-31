package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"forum/database"
	"forum/handlers"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = database.Init("storage/forum.db")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	defer db.Close()

	// Restrict arguments parsed
	if len(os.Args) != 1 {
		log.Println("Too many arguments")
		log.Println("Usage: go run main.go")
		return
	}

	// Candle hundler functions
	http.HandleFunc("/static/", handlers.StaticHandler)
	// http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/forgot-password", handlers.ForgotPasswordHandler)
	http.HandleFunc("/register", handlers.RegistrationHandler)

	// Inform user initialization of server
	log.Println("Server started on port 8080")

	// Start the server, handle emerging errors
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Failed to start server: ", err)
		return
	}
}
