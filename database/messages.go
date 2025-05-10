package database

import (
	"time"

	"forum/models"
)

// InsertMessage inserts a new message into the database.
func InsertMessage(senderID, receiverID int, content string, timestamp time.Time) error {
	_, err := db.Exec(`
        INSERT INTO messages (sender_id, receiver_id, content, timestamp)
        VALUES (?, ?, ?, ?)
    `,
		senderID,
		receiverID,
		content,
		timestamp,
	)
	return err
}

// GetLastMessage gets the last message between two users.
func GetLastMessage(userID1, userID2 int) (models.Message, error) {
	query := `
		SELECT
			m.id,
			m.sender_id,
			m.receiver_id,
			m.content,
			m.timestamp,
			u.username as sender_nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
			OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.timestamp DESC
		LIMIT 1
	`

	var message models.Message
	var timestampStr string // SQLite returns dates as strings

	err := db.QueryRow(
		query,
		userID1, userID2,
		userID2, userID1,
	).Scan(
		&message.ID,
		&message.SenderID,
		&message.ReceiverID,
		&message.Content,
		&timestampStr,
		&message.SenderNickname,
	)
	if err != nil {
		return models.Message{}, err
	}

	// Parse the timestamp
	message.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		// If parsing fails, just use current time
		message.Timestamp = time.Now()
	}

	return message, nil
}

// GetMessages retrieves messages between two users, ordered by timestamp.
// GetMessages retrieves messages between two users with proper sorting
func GetMessages(senderID, receiverID int, offset, limit int) ([]models.Message, error) {
	var messages []models.Message

	// Modify the SQL query to properly sort messages by timestamp in ascending order (oldest first)
	// Using SQLite's ? placeholders instead of PostgreSQL's $1, $2, etc.
	query := `
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.timestamp, u.username as sender_nickname
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.timestamp ASC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, senderID, receiverID, receiverID, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.Message
		var timestampStr string // SQLite returns dates as strings

		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &timestampStr, &msg.SenderNickname)
		if err != nil {
			return nil, err
		}

		// Parse the timestamp
		msg.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			// If parsing fails, just use current time
			msg.Timestamp = time.Now()
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
