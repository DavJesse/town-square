package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Message represents a chat message
type Message struct {
	Type       string    `json:"type"`
	Content    string    `json:"content,omitempty"`
	SenderID   int64     `json:"sender_id,omitempty"`
	ReceiverID int64     `json:"receiver_id,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	TempID     string    `json:"temp_id,omitempty"`
}

func (m *Message) Serialize() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Failed to serialize message: %v", err)
		return nil
	}
	return data
}
type Conversation struct {
	OtherUserID     int64     `json:"other_user_id"`
	Username        string    `json:"username"`
	IsOnline        bool      `json:"is_online"`
	LastSeen        time.Time `json:"last_seen"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
}

// Add a new message type for status updates
type StatusUpdateMessage struct {
	Type     string    `json:"type"`
	UserID   int64     `json:"user_id"`
	IsOnline bool      `json:"is_online"`
	LastSeen time.Time `json:"last_seen"`
}

// Add a new type for heartbeat messages
type HeartbeatMessage struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
