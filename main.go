package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/database"

	"forum/handlers/routes"
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

	// Configure all routes
	routes.RegisterRoutes()

	// start the server, handle emerging errors
	fmt.Printf("Server runing on http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Failed to start server: ", err)
		return
	}
}
