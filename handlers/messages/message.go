// real/chat/chat.go
package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"forum/database"
	errs "forum/handlers/errors"
	"forum/models"

	"github.com/gorilla/websocket"
)

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow all origins
}

// UserConnections holds the WebSocket connections for each user.
// key = userID
var (
	UserConnections = make(map[int]models.WSUser)
	userConnMutex   = &sync.Mutex{} // Mutex to protect concurrent access to UserConnections
)

// HandleWebSocket handles WebSocket connections for sending and receiving messages.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Get the user ID from the request context. This should have been added by the AuthMiddleware.
	sessionWithUsername, ok := r.Context().Value(database.SESSION_KEY).(*models.SessionWithUsername)
	if !ok {
		log.Println("User ID not found in request context")
		conn.Close()
		return
	}
	userID := sessionWithUsername.UserID

	// Store the connection.
	currentUser, err := database.GetUserByID(userID)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}

	// Send auth_success message immediately after connection
	authMessage := map[string]interface{}{
		"type":   "auth_success",
		"userID": userID,
	}

	if err := conn.WriteJSON(authMessage); err != nil {
		log.Printf("Error sending auth success message: %v", err)
		conn.Close()
		// delete(UserConnections, userID)
		// database.UpdateUserOnlineStatus(userID, false)
		return
	} else {
		userConnMutex.Lock()
		UserConnections[userID] = models.WSUser{
			User: currentUser,
			Conn: conn,
		}
		userConnMutex.Unlock()
		log.Printf("User %d connected\n", userID)
		database.UpdateUserOnlineStatus(userID, true)
		BroadcastUserStatus(userID, true)
	}

	defer func() {
		conn.Close()
		userConnMutex.Lock()
		delete(UserConnections, userID)
		userConnMutex.Unlock()
		database.UpdateUserOnlineStatus(userID, false)
		// Broadcast offline status to all other users
		BroadcastUserStatus(userID, false)
		log.Printf("User %d disconnected\n", userID)
	}()

	for {
		// Read message from the WebSocket connection.
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break // Exit the loop on error to close the connection.
		}

		// Parse the incoming message
		var msgData map[string]interface{}
		if err := json.Unmarshal(p, &msgData); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			continue
		}

		// Check message type
		msgType, ok := msgData["type"].(string)
		if !ok {
			log.Printf("Message type not found or not a string")
			continue
		}

		if msgType == "message" {
			// Extract the message data
			msgContent, ok := msgData["message"].(map[string]interface{})
			if !ok {
				log.Printf("Message content not found or not a map")
				continue
			}

			// Create the WSMessage from the extracted data
			var wsMessage models.WSMessage

			// Convert sender_id to int
			senderIDFloat, ok := msgContent["sender_id"].(float64)
			if !ok {
				log.Printf("Sender ID not found or not a number")
				continue
			}
			wsMessage.SenderID = int(senderIDFloat)

			// Convert receiver_id to int
			receiverIDFloat, ok := msgContent["receiver_id"].(float64)
			if !ok {
				log.Printf("Receiver ID not found or not a number")
				continue
			}
			wsMessage.ReceiverID = int(receiverIDFloat)

			// Get the content
			wsMessage.Content, ok = msgContent["content"].(string)
			if !ok {
				log.Printf("Content not found or not a string")
				continue
			}

			// Get current time
			timestamp := time.Now()
			wsMessage.Timestamp = timestamp.Format(time.RFC3339)

			// Get Sender Nickname from the database.
			user, err := database.GetUserByID(wsMessage.SenderID)
			if err != nil {
				log.Printf("Error getting Sender Nickname: %v", err)
				continue
			}
			wsMessage.SenderNickname = user.Nickname

			// Store the message in the database.
			err = database.InsertMessage(wsMessage.SenderID, wsMessage.ReceiverID, wsMessage.Content, timestamp)
			if err != nil {
				log.Printf("Error inserting message into database: %v", err)
				continue
			}

			// Format message for sending
			messageToSend := map[string]interface{}{
				"type":    "message",
				"message": wsMessage,
			}

			messageJSON, err := json.Marshal(messageToSend)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			// Send to receiver if online
			userConnMutex.Lock()
			receiverUser, ok := UserConnections[wsMessage.ReceiverID]
			userConnMutex.Unlock()
			if ok {
				if err := receiverUser.Conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
					log.Printf("Error sending message to receiver %d: %v", wsMessage.ReceiverID, err)
				}
			}

			// Send back to sender for confirmation
			if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				log.Printf("Error sending message back to sender: %v", err)
			}
		} else if msgType == "typing" || msgType == "stop_typing" {
			// Extract sender and receiver IDs
			senderIDFloat, ok := msgData["sender_id"].(float64)
			if !ok {
				log.Printf("Sender ID not found or not a number")
				continue
			}
			senderID := int(senderIDFloat)

			receiverIDFloat, ok := msgData["receiver_id"].(float64)
			if !ok {
				log.Printf("Receiver ID not found or not a number")
				continue
			}
			receiverID := int(receiverIDFloat)

			// Get sender's nickname
			sender, err := database.GetUserByID(senderID)
			if err != nil {
				log.Printf("Error getting sender info: %v", err)
				continue
			}

			// Add sender's nickname to the message
			msgData["sender_nickname"] = sender.Nickname

			// Forward the typing event to the receiver if they're online
			userConnMutex.Lock()
			receiverUser, ok := UserConnections[receiverID]
			userConnMutex.Unlock()

			if ok {
				// Send the typing event to the receiver
				if err := receiverUser.Conn.WriteJSON(msgData); err != nil {
					log.Printf("Error sending typing event to receiver %d: %v", receiverID, err)
				}
			}
		}
	}
}

// braodcast to all users when a status of a single user changes
func BroadcastUserStatus(userID int, isOnline bool) {
	// Create status change message
	statusMsg := map[string]interface{}{
		"type": "user_status_change",
		"user": map[string]interface{}{
			"id":        userID,
			"is_online": isOnline,
		},
	}

	messageJSON, err := json.Marshal(statusMsg)
	if err != nil {
		log.Printf("Error marshalling status message: %v", err)
		return
	}

	// Broadcast to all connected users
	userConnMutex.Lock()
	// Create a copy of the connections to avoid holding the lock during I/O
	connections := make(map[int]models.WSUser)
	for id, wsUser := range UserConnections {
		if id != userID {
			connections[id] = wsUser
		}
	}
	userConnMutex.Unlock()

	// Send messages without holding the lock
	for id, wsUser := range connections {
		if err := wsUser.Conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
			log.Printf("Error sending status update to user %d: %v", id, err)
		}
	}
}

// GetMessagesHandler retrieves the message history between two users.
func GetMessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the sender and receiver IDs from the request.
		sessionWithUsername, ok := r.Context().Value(database.SESSION_KEY).(*models.SessionWithUsername)
		if !ok {
			errs.JSONErrorResponse(w, "User ID not found in request context", http.StatusBadRequest)
			return
		}
		senderID := sessionWithUsername.UserID

		receiverID := r.URL.Query().Get("receiver_id")
		if receiverID == "" {
			errs.JSONErrorResponse(w, "Receiver ID is required", http.StatusBadRequest)
			return
		}
		// convert receiverID to int.
		var receiverIDInt int
		fmt.Sscanf(receiverID, "%d", &receiverIDInt)

		// Get lastMessageID and limit from query parameters
		lastMessageIDStr := r.URL.Query().Get("last_message_id")
		limitStr := r.URL.Query().Get("limit")
		direction := r.URL.Query().Get("direction")

		// Default values
		lastMessageID := 0
		limit := 10

		// Parse lastMessageID if provided
		if lastMessageIDStr != "" {
			fmt.Sscanf(lastMessageIDStr, "%d", &lastMessageID)
		}

		// Parse limit if provided
		if limitStr != "" {
			fmt.Sscanf(limitStr, "%d", &limit)
		}

		// Default direction is "older" if not specified
		if direction == "" {
			direction = "older"
		}

		messages, err := database.GetMessages(senderID, receiverIDInt, lastMessageID, limit, direction)
		if err != nil {
			log.Printf("Error getting messages: %v", err)
			errs.JSONErrorResponse(w, "Failed to get messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

// GetAllUsersHandler retrieves a list of all users with their online status and chat history info
// - Online/offline status based on last_active and active WebSocket connections
// - Chat history information for message sorting
// - Last message data for each user with chat history
func GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionWithUsername, ok := r.Context().Value(database.SESSION_KEY).(*models.SessionWithUsername)
		if !ok {
			errs.JSONErrorResponse(w, "User ID not found in request context", http.StatusBadRequest)
			return
		}
		userID := sessionWithUsername.UserID

		users, err := database.GetAllUsers(userID)
		if err != nil {
			log.Printf("Error getting users: %v", err)
			errs.JSONErrorResponse(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		// Update online status based on active WebSocket connections
		// This ensures real-time online status tracking beyond just the database timestamp
		userConnMutex.Lock()
		for index, user := range users {
			if _, ok = UserConnections[user.ID]; ok {
				users[index].IsOnline = true
				// Also update the database status to match
				database.UpdateUserOnlineStatus(user.ID, true)
			}
		}
		userConnMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func SendMessageHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var msgData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&msgData); err != nil {
			errs.JSONErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID from context
		sessionWithUsername, ok := r.Context().Value(database.SESSION_KEY).(*models.SessionWithUsername)
		if !ok {
			errs.JSONErrorResponse(w, "User ID not found in request context", http.StatusBadRequest)
			return
		}
		senderID := sessionWithUsername.UserID

		// Extract the message data
		msgType, ok := msgData["type"].(string)
		if !ok || msgType != "message" {
			errs.JSONErrorResponse(w, "Invalid message type", http.StatusBadRequest)
			return
		}

		msgContent, ok := msgData["message"].(map[string]interface{})
		if !ok {
			errs.JSONErrorResponse(w, "Invalid message content", http.StatusBadRequest)
			return
		}

		// Validate the message content
		receiverIDFloat, ok := msgContent["receiver_id"].(float64)
		if !ok {
			errs.JSONErrorResponse(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		receiverID := int(receiverIDFloat)

		content, ok := msgContent["content"].(string)
		if !ok || content == "" {
			errs.JSONErrorResponse(w, "Message content cannot be empty", http.StatusBadRequest)
			return
		}

		// Get current time
		timestamp := time.Now()

		// Get Sender Nickname from the database
		user, err := database.GetUserByID(senderID)
		if err != nil {
			log.Printf("Error getting Sender Nickname: %v", err)
			errs.JSONErrorResponse(w, "Failed to get sender info", http.StatusInternalServerError)
			return
		}

		// Create the message
		wsMessage := models.WSMessage{
			SenderID:       senderID,
			ReceiverID:     receiverID,
			Content:        content,
			Timestamp:      timestamp.Format(time.RFC3339),
			SenderNickname: user.Nickname,
		}

		// Store the message in the database
		err = database.InsertMessage(wsMessage.SenderID, wsMessage.ReceiverID, wsMessage.Content, timestamp)
		if err != nil {
			log.Printf("Error inserting message into database: %v", err)
			errs.JSONErrorResponse(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		// Format message for sending
		messageToSend := map[string]interface{}{
			"type":    "message",
			"message": wsMessage,
		}

		// Try to send to receiver if online
		userConnMutex.Lock()
		receiverUser, online := UserConnections[receiverID]
		userConnMutex.Unlock()

		if online {
			messageJSON, _ := json.Marshal(messageToSend)
			if err := receiverUser.Conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				log.Printf("Error sending message to receiver %d: %v", receiverID, err)
				// Continue anyway, as the message is stored in the database
			}
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(messageToSend)
	}
}

func LogoutWS(userID int) {
	// If user has an active websocket connection, close it
	userConnMutex.Lock()
	wsUser, exists := UserConnections[userID]
	if exists {
		// Get the connection but delete from map under lock
		delete(UserConnections, userID)
		userConnMutex.Unlock()

		// Close the connection outside the lock
		wsUser.Conn.Close()
	} else {
		userConnMutex.Unlock()
	}

	// Broadcast status change to all connected users
	BroadcastUserStatus(userID, false)
}
