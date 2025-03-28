package database

import (
	"fmt"
	"log"

	"forum/models"
)

// SendPublicMessage sends a public message to every online user
func SendPublicMessage(senderID int, message string) error {
	// Get the list of online users
	users, err := GetOnlineUsers()
	if err != nil {
		log.Println("Error fetching online users:", err)
		return err
	}

	// Iterate over the online users and send the message to each of them
	for _, user := range users {
		// Skip the sender
		if user.ID == senderID {
			continue
		}

		// Insert the message into the database for each user
		err := SendPrivateMessage(senderID, user.ID, message)
		if err != nil {
			log.Printf("Error sending message to user %d: %v\n", user.ID, err)
			// Optionally, you can continue the loop even if there's an error for a particular user
			continue
		}
	}

	return nil
}

// SendMessage sends a private message from one user to another
func SendPrivateMessage(senderID, receiverID int, message string) error {
	fmt.Printf("SENDER: %d, RECEIVER: %d, MESSAGE: %v\n", senderID, receiverID, message)
	query := `INSERT INTO messages (sender_id, receiver_id, message) VALUES (?, ?, ?)`
	_, err := db.Exec(query, senderID, receiverID, message)
	return err
}

// GetMessages retrieves all messages between two users.
func GetMessages(user1ID, user2ID int, limit, offset int) ([]models.Message, error) {
	query := `SELECT id, sender_id, receiver_id, message, created_at
	          FROM messages 
	          WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
	          ORDER BY created_at DESC
	          LIMIT ? OFFSET ?`

	rows, err := db.Query(query, user1ID, user2ID, user2ID, user1ID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// SetUserOnlineStatus updates the user's online status.
func SetUserOnlineStatus(userID int64, isOnline bool) error {
	query := `INSERT INTO user_status (user_id, is_online, last_seen) 
              VALUES (?, ?, CURRENT_TIMESTAMP) 
              ON CONFLICT(user_id) 
              DO UPDATE SET is_online = ?, last_seen = CURRENT_TIMESTAMP`
	_, err := db.Exec(query, userID, isOnline, isOnline)
	return err
}

// GetOnlineUsers retrieves all online users.
func GetOnlineUsers() ([]models.User, error) {
	query := `SELECT u.id, u.username, u.first_name, u.last_name, u.image 
	          FROM users u 
	          JOIN user_status us ON u.id = us.user_id 
	          WHERE us.is_online = TRUE`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Image); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Update the user's status when they disconnect
func UpdateUserOfflineStatus(userID int) {
	err := SetUserOnlineStatus(int64(userID), false)
	if err != nil {
		log.Println("Error setting user status to offline:", err)
	}
}
