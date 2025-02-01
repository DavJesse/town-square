package models

import (
	"time"
)

type WebError struct {
	Code  int
	Issue string
}

// user struct
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Exclude from JSON responses
	Bio       *string   `json:"bio,omitempty"`
	Image     *string   `json:"image,omitempty"`
	SessionID *string   `json:"session_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	Expiry    time.Time `json:"expiry"`
}

// category struct
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Post struct
type Post struct {
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Media     string    `json:"media,omitempty"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// PostCategory struct
type PostCategory struct {
	PostID     int `json:"post_id"`
	CategoryID int `json:"category_id"`
}

// Comment Struct
type Comment struct {
	UUID      string    `json:"uuid"`
	Content   string    `json:"content"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Like struct
type Like struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	PostID    *int `json:"post_id,omitempty"`
	CommentID *int `json:"comment_id,omitempty"`
}

// Dislike struct
type Dislike struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	PostID    *int `json:"post_id,omitempty"`
	CommentID *int `json:"comment_id,omitempty"`
}
