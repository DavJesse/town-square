package models

import "time"

// Message represents a private message between two users.
type Message struct {
	ID         int64     `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

// UserStatus represents the online status of a user.
type UserStatus struct {
	UserID   int64     `json:"user_id"`
	IsOnline bool      `json:"is_online"`
	LastSeen time.Time `json:"last_seen"`
}

type WebSocketOnConnectionData struct {
	MyID        int    `json:"userid"`
	MyName      string    `json:"username"`
	OnlineUsers []User `json:"online_users"`
}
