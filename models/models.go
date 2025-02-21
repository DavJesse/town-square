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
	Bio       string    `json:"bio,omitempty"`
	Image     string    `json:"image,omitempty"`
	SessionID string    `json:"session_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	Expiry    time.Time `json:"expiry"`
}
type SessionWithUsername struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
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

// used to fetch the post with the creator of the post from the database
type PostWithUsername struct {
	UUID          string    `json:"uuid"`
	Creator       string    `json:"creator"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Media         string    `json:"media,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	LikesCount    int       `json:"likes_count"`
	DislikesCount int       `json:"dislikes_count"`
}

// PostWithCategories struct to hold post data along with categories, likes, dislikes, and comments
type PostWithCategories struct {
	UUID          string               `json:"uuid"`
	Title         string               `json:"title"`
	Content       string               `json:"content"`
	Media         string               `json:"media"`
	Username      string               `json:"username"`
	UserID        int                  `json:"user_id"`
	CreatedAt     time.Time            `json:"created_at"`
	Categories    []string             `json:"categories"` // Categories as a slice of strings
	LikesCount    int                  `json:"likes_count"`
	DislikesCount int                  `json:"dislikes_count"`
	Comments      []CommentWithCreator `json:"comments"` // Slice of comments
}

// Comment struct to hold comment data
type Comment struct {
	UUID      string    `json:"uuid"`
	Content   string    `json:"content"`
	PostID    string    `json:"post_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentWithCreator struct {
	UUID          string    `json:"uuid"`
	Content       string    `json:"content"`
	PostID        string    `json:"post_id"`
	Creator       string    `json:"creator"`
	CreatedAt     time.Time `json:"created_at"`
	LikesCount    int       `json:"likes_count"`
	DislikesCount int       `json:"dislikes_count"`
}
