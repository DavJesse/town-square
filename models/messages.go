package models

import (
	"github.com/gorilla/websocket"
	"time"
)

// Message struct to hold message data
type Message struct {
	ID             int       `json:"id"`
	SenderID       int       `json:"sender_id"`
	ReceiverID     int       `json:"receiver_id"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
	SenderNickname string    `json:"sender_nickname"`
}

// WSUser holds a user and their WebSocket connection
type WSUser struct {
	User User
	Conn *websocket.Conn
}

// Define a struct to hold validation errors.  This will make it easier to
// send structured error responses.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Message struct to hold message data for WebSocket communication
type WSMessage struct {
	SenderID       int    `json:"sender_id"`
	ReceiverID     int    `json:"receiver_id"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"` // Store as string for easy transfer
	SenderNickname string `json:"sender_nickname"`
}
