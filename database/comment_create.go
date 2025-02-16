package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func CreateComment(userID int, postID string, content string) (int64, error) {
	uuid, err := GenerateUUID()
	if err != nil {
		return 0, err
	}
	query := `INSERT INTO comments (uuid, user_id, post_id, content) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, uuid, userID, postID, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GenerateUUID() (string, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating UUID: %w", err)
	}
	return newUUID.String(), nil
}
