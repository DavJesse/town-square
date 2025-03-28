package messages

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"forum/database"
	"forum/models"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

var (
	server   *WebSocketServer
	upgrader websocket.Upgrader
)

// Constructor for creating a new WebSocketServer instance
func NewWebSocketServer() {
	server = &WebSocketServer{
		clients: make(map[*websocket.Conn]bool),
	}
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}

// WebSocketHandler handles new WebSocket connections and message broadcasting
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the session
	// sessionID, username, err := database.GetUserData(r)
	userID, userName, err := database.GetUserData(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Locking to safely add the new client to the clients map
	server.mu.Lock()
	server.clients[conn] = true
	server.mu.Unlock()

	errSetOnlineStatus := database.SetUserOnlineStatus(int64(userID), true)
	if errSetOnlineStatus != nil {
		fmt.Printf("[ERROR] Failed to update user online status")
	} else {
		fmt.Printf("Username: %s, UserID: %d\n", userName, userID)
	}

	// Send the list of online users to the newly connected user
	onlineUsers, err := database.GetOnlineUsers()
	if err != nil {
		log.Println("Error retrieving online users:", err)
	}

	fmt.Println("ONLINE: ", onlineUsers)
	// Send the online users list to the connected client
	err = conn.WriteJSON(onlineUsers)
	if err != nil {
		log.Println("Error sending online users to client:", err)
		conn.Close()
		return
	}

	// Listen for incoming messages
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Handle incoming message (store in DB and broadcast)
		err = database.SendMessage(message.SenderID, message.ReceiverID, message.Message)
		if err != nil {
			log.Println("Error storing message:", err)
			continue
		}

		// Lock to safely broadcast message to all clients
		server.mu.Lock()
		for client := range server.clients {
			if err := client.WriteJSON(message); err != nil {
				log.Println("Error sending message to client:", err)
				client.Close()
				delete(server.clients, client)
			}
		}
		server.mu.Unlock()
	}
}
