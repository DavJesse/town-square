package models

import "time"

// Message struct to hold message data
type Message struct {
	ID             int       `json:"id"`
	SenderID       int       `json:"sender_id"`
	ReceiverID     int       `json:"receiver_id"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
	SenderNickname string    `json:"sender_nickname"`
}
