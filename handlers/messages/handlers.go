package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/database"
	"forum/models"
)

func GetConversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get userID using the ServeWS pattern
	userID, err := GetUserIDFromRequest(r)
	fmt.Println("UserID: ", userID)
	if err != nil {
		fmt.Printf("Failed to get userID: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"error":         "Authentication required",
			"status":        "error",
			"conversations": []any{},
		})
		return
	}

	fmt.Printf("Getting conversations for user: %d", userID)

	conversations, err := database.GetConversations(userID)
	if err != nil {
		fmt.Printf("Failed to get conversations: %v", err)
		fmt.Printf("Stack trace: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error":         fmt.Sprintf("Failed to get conversations: %v", err),
			"status":        "error",
			"conversations": []any{},
		})
		return
	}

	if conversations == nil {
		fmt.Printf("No conversations found, returning empty array")
		conversations = []models.Conversation{}
	}

	response := map[string]any{
		"status":        "success",
		"conversations": conversations,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to encode response: %v", err)
		return
	}

	fmt.Printf("=== Completed GetConversationsHandler successfully ===")
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("=== Starting GetMessagesHandler ===")
	fmt.Printf("Request Method: %s, URL: %s", r.Method, r.URL.Path)

	// Get userID using the ServeWS pattern
	userID, err := GetUserIDFromRequest(r)
	if err != nil {
		fmt.Printf("Failed to get userID: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	fmt.Printf("URL parts: %v", parts)

	if len(parts) < 4 {
		fmt.Printf("Invalid URL path length: %d, parts: %v", len(parts), parts)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	otherUserID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		fmt.Printf("Failed to parse otherUserID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("Retrieving messages between users %d and %d", userID, otherUserID)

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
			fmt.Printf("Using page number: %d", page)
		} else {
			fmt.Printf("Invalid page parameter: %s", pageStr)
		}
	}

	messages, err := database.GetMessages(userID, otherUserID, page)
	if err != nil {
		fmt.Printf("Failed to get messages: %v", err)
		fmt.Printf("Stack trace: %+v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Retrieved %d messages successfully", len(messages))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		fmt.Printf("Failed to encode messages: %v", err)
		return
	}

	fmt.Printf("=== Completed GetMessagesHandler successfully ===")
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("=== Starting SendMessageHandler ===")
	fmt.Printf("Request Method: %s, URL: %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userID using the ServeWS pattern
	senderID, err := GetUserIDFromRequest(r)
	if err != nil {
		fmt.Printf("Failed to get userID: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the receiver ID from the URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	receiverID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse the message content from the request body
	var messageData struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&messageData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if messageData.Content == "" {
		http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
		return
	}

	// Create and store the message
	message := &models.Message{
		Type:       "message",
		Content:    messageData.Content,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Timestamp:  time.Now().UTC(),
	}

	if err := database.StoreMessage(message); err != nil {
		fmt.Printf("Failed to store message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Broadcast the message to connected clients
	if hub != nil {
		hub.Broadcast <- message
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": message,
	})

	fmt.Printf("=== Completed SendMessageHandler successfully ===")
}
