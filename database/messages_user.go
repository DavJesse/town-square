package database

import (
	"fmt"
	"sort"
	"time"

	"forum/models"
)

// InsertUser inserts a new user into the database.
func InsertUser(nickname string, age int, gender string, firstName string, lastName string, email string, passwordHash string) error {
	_, err := db.Exec(`
        INSERT INTO users (nickname, age, gender, first_name, last_name, email, password_hash, last_active, is_online)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
		nickname,
		age,
		gender,
		firstName,
		lastName,
		email,
		passwordHash,
		time.Now(),
		false,
	)
	return err
}

// GetUserByEmail retrieves a user from the database by email.
func GetUserByEmail(email string) (models.UserWS, error) {
	var user models.UserWS
	err := db.QueryRow("SELECT id, username, password_hash FROM users WHERE email = ?", email).Scan(&user.ID, &user.Nickname, &user.PasswordHash)
	return user, err
}

// GetUserByID retrieves a user from the database by ID.
func GetUserByID(id int) (models.UserWS, error) {
	fmt.Println("GETUSERBYID: ", id)
	var user models.UserWS
	err := db.QueryRow("SELECT id, username, first_name, last_name, email, gender, age FROM users WHERE id = ?", id).Scan(&user.ID, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Gender, &user.Age)
	fmt.Println("FETCHED_USER: ", user);
	return user, err
}

// GetUserIDByToken retrieves a user ID from the database, given a session token.
func GetUserIDByToken(token string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userID)
	return userID, err
}

// GetAllUsers: retrieves the online/offline users.
//
//	returns all users except the current user
//
// GetAllUsers retrieves all users with better classification and sorting information.
//
// Returns all users except the current user, with:
// - Online/offline status
// - Last message information for sorting
// - Indication of whether you've exchanged messages with them
func GetAllUsers(currentUserID int) ([]models.UserWS, error) {
	// SQLite compatible query to get all users except the current user, with their online status
	query := `
		SELECT 
			u.id, 
			u.username,
			u.first_name,
			u.last_name, 
			u.email,
			u.gender,
			u.age,
			CASE WHEN datetime(u.last_active) > datetime('now', '-5 minutes') THEN 1 ELSE 0 END as is_online,
			u.last_active as last_seen,
			CASE WHEN EXISTS (
				SELECT 1 FROM messages m 
				WHERE (m.sender_id = u.id AND m.receiver_id = ?) 
				OR (m.sender_id = ? AND m.receiver_id = u.id)
			) THEN 1 ELSE 0 END as has_messages
		FROM users u
		WHERE u.id != ?
		ORDER BY 
			has_messages DESC, -- First show users with message history
			nickname ASC      -- Then alphabetically by nickname
	`

	rows, err := db.Query(query, currentUserID, currentUserID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserWS
	for rows.Next() {
		var user models.UserWS
		var isOnlineInt int
		var hasMessagesInt int
		var lastSeenStr string // SQLite returns dates as strings

		err := rows.Scan(
			&user.ID,
			&user.Nickname,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Gender,
			&user.Age,
			&isOnlineInt,
			&lastSeenStr,
			&hasMessagesInt,
		)
		if err != nil {
			return nil, err
		}

		user.IsOnline = isOnlineInt == 1
		user.HasChatHistory = hasMessagesInt == 1

		// Parse the last_seen timestamp
		if lastSeenStr != "" {
			user.LastSeen, err = time.Parse("2006-01-02 15:04:05", lastSeenStr)
			if err != nil {
				// If parsing fails, just use current time
				user.LastSeen = time.Now()
			}
		}

		users = append(users, user)
	}

	// For each user, get their last message with the current user
	for i, user := range users {
		fmt.Println("USER: ", user.FirstName)
		if user.HasChatHistory {
			lastMessage, err := GetLastMessage(currentUserID, user.ID)
			if err == nil && lastMessage.ID > 0 {
				users[i].LastMessage = lastMessage
			}
		}
	}

	// Sort users properly: first by message existence and timestamp, then alphabetically
	// Note: We need to sort in-memory to properly handle the LastMessage timestamp
	sort.Slice(users, func(i, j int) bool {
		// If one has messages and the other doesn't, prioritize the one with messages
		if users[i].HasChatHistory != users[j].HasChatHistory {
			return users[i].HasChatHistory
		}

		// If both have message history, sort by most recent message
		if users[i].HasChatHistory && users[j].HasChatHistory {
			// Check if LastMessage is properly loaded for both
			if users[i].LastMessage.ID > 0 && users[j].LastMessage.ID > 0 {
				return users[i].LastMessage.Timestamp.After(users[j].LastMessage.Timestamp)
			}
		}

		// Default: sort alphabetically by nickname
		return users[i].Nickname < users[j].Nickname
	})

	return users, nil
}

// UpdateUserOnlineStatus updates the last_active timestamp for a user to mark them as online
func UpdateUserOnlineStatus(userID int, isOnline bool) error {
	var err error

	if isOnline {
		// If marking user as online, update last_active to current time
		_, err = db.Exec("UPDATE users SET last_active = ? AND is_online = ? WHERE id = ?", time.Now(), isOnline, userID)
	} else {
		// If explicitly marking user as offline, set last_active to a time in the past
		// This ensures they'll be shown as offline immediately
		pastTime := time.Now().Add(-10 * time.Minute)
		_, err = db.Exec("UPDATE users SET last_active = ? AND is_online = ? WHERE id = ?", pastTime, isOnline, userID)
	}

	return err
}
