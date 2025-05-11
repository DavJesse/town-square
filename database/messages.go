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

// GetMessages retrieves messages between two users using ID-based pagination
func GetMessages(senderID, receiverID int, lastMessageID int, limit int, direction string) ([]models.Message, error) {
	var messages []models.Message
	var query string
	var args []interface{}

	// Base condition for conversation between the two users
	baseCondition := "(m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)"
	baseArgs := []interface{}{senderID, receiverID, receiverID, senderID}

	if lastMessageID > 0 && direction == "older" {
		// Get messages older than the lastMessageID (for scrolling up)
		query = `
			SELECT m.id, m.sender_id, m.receiver_id, m.content, m.timestamp, u.username as sender_nickname
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE (` + baseCondition + `) AND m.id < ?
			ORDER BY m.id ASC
			LIMIT ?
		`
		args = append(baseArgs, lastMessageID, limit)
	} else {
		// Initial load - get the most recent messages
		query = `
			SELECT m.id, m.sender_id, m.receiver_id, m.content, m.timestamp, u.username as sender_nickname
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE ` + baseCondition + `
			ORDER BY m.id DESC
			LIMIT ?
		`
		args = append(baseArgs, limit)
	}

	rows, err := db.Query(query, args...)
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

	// If this is the initial load, we need to reverse the order
	// since we queried with ORDER BY id DESC
	if lastMessageID == 0 {
		// Reverse the slice
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
	}

	return messages, nil
}
