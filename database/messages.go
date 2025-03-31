package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/models"
)

// StoreMessage stores a message in the database
func StoreMessage(message *models.Message) error {
	_, err := db.Exec(`
        INSERT INTO messages (sender_id, receiver_id, content, created_at)
        VALUES (?, ?, ?, ?)
    `, message.SenderID, message.ReceiverID, message.Content, message.Timestamp)
	return err
}

// UpdateUserStatus updates a user's online status and last seen time
func UpdateUserStatus(userID int64, isOnline bool) error {
	currentTime := time.Now().UTC()

	// First try to update existing record
	result, err := db.Exec(`
        UPDATE user_status 
        SET is_online = ?,
            last_seen = CASE 
                WHEN is_online = 0 OR is_online <> ? THEN ?
                ELSE last_seen 
            END
        WHERE user_id = ?
    `, isOnline, isOnline, currentTime, userID)
	if err != nil {
		fmt.Printf("Failed to update user status: %v", err)
		return err
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// No existing record, insert a new one
		_, err = db.Exec(`
            INSERT INTO user_status (user_id, is_online, last_seen)
            VALUES (?, ?, ?)
        `, userID, isOnline, currentTime)
		if err != nil {
			fmt.Printf("Failed to insert user status: %v", err)
			return err
		}
	}

	return nil
}

// GetUserLastSeen retrieves a user's last seen time from the database
func GetUserLastSeen(userID int64) (time.Time, error) {
	var lastSeen time.Time
	err := db.QueryRow("SELECT last_seen FROM user_status WHERE user_id = ?", userID).Scan(&lastSeen)
	if err != nil {
		return time.Now().UTC(), err
	}
	return lastSeen, nil
}

// Define Message and Conversation structs

type Message struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// GetMessages retrieves messages between two users

func GetMessages(userID, otherUserID int64, page int) ([]Message, error) {
	offset := (page - 1) * 10
	rows, err := db.Query(`
        SELECT id, sender_id, receiver_id, content, created_at 
        FROM messages 
        WHERE (sender_id = ? AND receiver_id = ?) 
           OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at DESC
        LIMIT 10 OFFSET ?
    `, userID, otherUserID, otherUserID, userID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// GetConversations retrieves conversations for a user

func GetConversations(userID int64) ([]models.Conversation, error) {
	query := `
        SELECT DISTINCT 
            u.id as user_id,
            u.username,
            COALESCE(us.is_online, false) as is_online,
            COALESCE(strftime('%Y-%m-%d %H:%M:%f', us.last_seen), CURRENT_TIMESTAMP) as last_seen,
            (
                SELECT content 
                FROM messages m2 
                WHERE (m2.sender_id = ? AND m2.receiver_id = u.id) 
                   OR (m2.sender_id = u.id AND m2.receiver_id = ?)
                ORDER BY m2.created_at DESC 
                LIMIT 1
            ) as last_message,
            (
                SELECT created_at 
                FROM messages m3 
                WHERE (m3.sender_id = ? AND m3.receiver_id = u.id) 
                   OR (m3.sender_id = u.id AND m3.receiver_id = ?)
                ORDER BY m3.created_at DESC 
                LIMIT 1
            ) as last_message_time
        FROM messages m
        JOIN users u ON (
            CASE 
                WHEN m.sender_id = ? THEN m.receiver_id
                ELSE m.sender_id 
            END = u.id
        )
        LEFT JOIN user_status us ON us.user_id = u.id
        WHERE m.sender_id = ? OR m.receiver_id = ?
        GROUP BY u.id
        ORDER BY last_message_time DESC NULLS LAST
	`

	rows, err := db.Query(query, userID, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %v", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		var lastMessage sql.NullString
		var lastMessageTime sql.NullTime
		var lastSeenStr string // temporary variable to capture the string value

		err := rows.Scan(
			&conv.OtherUserID,
			&conv.Username,
			&conv.IsOnline,
			&lastSeenStr, // scan into string instead of time.Time
			&lastMessage,
			&lastMessageTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %v", err)
		}

		// Parse the last_seen string into a time.Time value.
		parsedTime, err := time.Parse("2006-01-02 15:04:05.000000", lastSeenStr)
		if err != nil {
			// If the layout doesn't match, try a simpler layout or default to time.Now()
			parsedTime, err = time.Parse("2006-01-02 15:04:05", lastSeenStr)
			if err != nil {
				parsedTime = time.Now()
			}
		}
		conv.LastSeen = parsedTime

		if lastMessage.Valid {
			conv.LastMessage = lastMessage.String
		}
		if lastMessageTime.Valid {
			conv.LastMessageTime = lastMessageTime.Time
		} else {
			conv.LastMessageTime = time.Now()
		}

		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversations: %v", err)
	}

	return conversations, nil
}

// GetUsersHandler returns a list of registered users

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Query the database for users
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		fmt.Printf("Failed to query users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build the list of users
	users := []map[string]any{}
	for rows.Next() {
		var id int
		var nickname string
		if err := rows.Scan(&id, &nickname); err != nil {
			fmt.Printf("Failed to scan user: %v", err)
			continue
		}
		users = append(users, map[string]any{
			"id":       id,
			"nickname": nickname,
		})
	}

	// Return the list of users as JSON
	json.NewEncoder(w).Encode(map[string]any{
		"users": users,
	})
}

// GetUserById retrieves a user by ID

func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userId := parts[len(parts)-1]

	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user struct {
		ID       int64  `json:"id"`
		Nickname string `json:"nickname"`
		IsOnline bool   `json:"is_online"`
		LastSeen string `json:"last_seen"`
	}

	err = db.QueryRow(`
		SELECT u.id, u.username, 
			   COALESCE(us.is_online, false) as is_online,
			   CASE 
				   WHEN us.is_online = true THEN datetime('now')
				   ELSE COALESCE(datetime(us.last_seen), datetime('now'))
			   END as last_seen
		FROM users u
		LEFT JOIN user_status us ON us.user_id = u.id
		WHERE u.id = ?`, id).Scan(&user.ID, &user.Nickname, &user.IsOnline, &user.LastSeen)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		fmt.Printf("Failed to fetch user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
