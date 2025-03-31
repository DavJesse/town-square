package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"forum/database"
	"forum/models"

	"github.com/gorilla/websocket"
)

type MessageHub struct {
	Clients    map[*Client]bool
	Broadcast  chan *models.Message
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.RWMutex
}

type Client struct {
	Hub      *MessageHub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   int64
	IsOnline bool
	LastSeen time.Time
}

var hub *MessageHub = nil // singleton

func NewMessageHub() {
	hub = &MessageHub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *models.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	// Start the hub's processing goroutine
	go hub.Run()
}

func (h *MessageHub) Run() {
	fmt.Printf("WebSocket MessageHub is running...")

	// Start a ticker to periodically check for inactive users
	statusCheckTicker := time.NewTicker(60 * time.Second)

	go func() {
		for range statusCheckTicker.C {
			h.checkInactiveUsers()
		}
	}()

	for {
		select {
		case client := <-h.Register:
			fmt.Printf("Registering new WebSocket client: %d", client.UserID)
			client.LastSeen = time.Now().UTC()
			h.Clients[client] = true
			// Update user status to online
			database.UpdateUserStatus(client.UserID, true)

		case client := <-h.Unregister:
			fmt.Printf("Unregistering WebSocket client: %d", client.UserID)
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				// Update user status to offline
				database.UpdateUserStatus(client.UserID, false)
			}

		case message := <-h.Broadcast:
			fmt.Printf("Broadcasting WebSocket message: %+v", message)
			h.handleMessage(message)
		}
	}
}

// ServeWS handles WebSocket requests from clients
func ServeWS(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session
	session, loggedIn := database.IsLoggedIn(r)

	var userID int64 = 0 // Default to 0 if not logged in
	var err error
	var userData models.User
	// Retrieve user data from session if logged in
	if loggedIn {
		userData, err = database.GetUserbySessionID(session.SessionID)
		if err != nil {
			fmt.Printf("Failed to get user data: %v", err)
			// Still attempt to upgrade, but with userID = 0
		} else {
			userID = int64(userData.ID)
		}
	} else {
		// User is not logged in, send a Javascript to redirect instead of StatusUnauthorized
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "window.location.href = '/login';")
		return
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	fmt.Printf("WebSocket connection established with userID: %d", userID)

	// Create new client using the properly typed userID
	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		IsOnline: true,
	}

	// Register client with hub
	hub.Register <- client

	// Start client routines
	go client.WritePump()
	go client.ReadPump()

	// Send initial message with userID after a short delay to ensure client is ready
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure connection is ready
		initialMsg := map[string]interface{}{
			"type":   "init",
			"userID": userID,
			"user": userData,
		}
		initialMsgJSON, err := json.Marshal(initialMsg)
		if err != nil {
			fmt.Printf("Error marshaling initial message: %v", err)
		} else {
			client.Send <- initialMsgJSON
		}
	}()
}

// GetUserIDFromRequest extracts user ID from the request using the same method as ServeWS
func GetUserIDFromRequest(r *http.Request) (int64, error) {
	session, loggedIn := database.IsLoggedIn(r)
	if !loggedIn {
		return 0, fmt.Errorf("unauthorized")
	}

	userData, err := database.GetUserbySessionID(session.SessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user data: %v", err)
	}

	return int64(userData.ID), nil
}

// Check for inactive clients and mark them offline
func (h *MessageHub) checkInactiveUsers() {
	threshold := time.Now().Add(-2 * time.Minute)

	for client := range h.Clients {
		if client.LastSeen.Before(threshold) {
			fmt.Printf("Client %d inactive, marking as offline", client.UserID)
			database.UpdateUserStatus(client.UserID, false)
		}
	}
}

func (h *MessageHub) handleMessage(message *models.Message) {
	// Handle heartbeat messages
	if message.Type == "heartbeat" {
		// Update the last seen time for the client
		for client := range h.Clients {
			if client.UserID == message.SenderID {
				client.LastSeen = time.Now().UTC()
				database.UpdateUserStatus(client.UserID, true)
				break
			}
		}
		return
	}

	// For regular messages, store in database
	if message.Type == "message" {
		if err := database.StoreMessage(message); err != nil {
			fmt.Printf("Failed to store message: %v", err)
			return
		}

		// Find receiver's client connection
		for client := range h.Clients {
			if client.UserID == message.ReceiverID || client.UserID == message.SenderID {
				select {
				case client.Send <- message.Serialize():
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// WritePump pumps messages from the Send channel to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(45 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The channel was closed.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			fmt.Printf("Sent message to client: %s\n", string(message))

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub's Broadcast channel.
func (c *Client) ReadPump() {
	defer func() {
		// Unregister the client and close the connection.
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512KB max message size
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		// Update last activity time on successful pong
		c.LastSeen = time.Now().UTC()
		database.UpdateUserStatus(c.UserID, true)
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var message models.Message
		if err := json.Unmarshal(data, &message); err != nil {
			fmt.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Set sender information and timestamp.
		message.SenderID = c.UserID
		if message.Timestamp.IsZero() {
			message.Timestamp = time.Now()
		}

		// Update the client's last seen time
		c.LastSeen = time.Now().UTC()

		c.Hub.Broadcast <- &message
	}
}

// Add this method to MessageHub
func (h *MessageHub) broadcastStatusUpdate(userID int64, isOnline bool) {
	// Get the current time for last_seen
	currentTime := time.Now().UTC()

	// Fetch the actual last seen time from the database for accuracy
	lastSeen, err := database.GetUserLastSeen(userID)
	if err != nil {
		fmt.Printf("Failed to get last seen time: %v", err)
		lastSeen = currentTime // Fallback
	}

	// Create status update message
	statusMsg := &models.StatusUpdateMessage{
		Type:     "status_update",
		UserID:   userID,
		IsOnline: isOnline,
		LastSeen: lastSeen,
	}

	// Serialize the status message
	msgData, err := json.Marshal(statusMsg)
	if err != nil {
		fmt.Printf("Failed to serialize status update: %v", err)
		return
	}

	// Send to all connected clients
	for client := range h.Clients {
		select {
		case client.Send <- msgData:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

// GetConversations retrieves conversations for a user using the ServeWS pattern
func GetConversations(w http.ResponseWriter, r *http.Request) {
	// Get user ID using the same method as ServeWS
	userID, err := GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conversations, err := database.GetConversations(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get conversations: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"conversations": conversations,
	})
}

// GetMessages retrieves messages between two users using the ServeWS pattern
func GetMessages(w http.ResponseWriter, r *http.Request, otherUserID int64, page int) {
	// Get user ID using the same method as ServeWS
	userID, err := GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	messages, err := database.GetMessages(userID, otherUserID, page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get messages: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
	})
}

// SendMessage handles sending a message using the ServeWS pattern
func SendMessage(w http.ResponseWriter, r *http.Request, receiverID int64) {
	// Get user ID using the same method as ServeWS
	senderID, err := GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var messageData struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&messageData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message := &models.Message{
		Type:       "message",
		Content:    messageData.Content,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Timestamp:  time.Now(),
	}

	if err := database.StoreMessage(message); err != nil {
		http.Error(w, fmt.Sprintf("Failed to store message: %v", err), http.StatusInternalServerError)
		return
	}

	// Broadcast the message to connected clients
	if hub != nil {
		hub.Broadcast <- message
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": message,
	})
}
