package database

import "forum/utils"

func CreateComment(userID int, postID string, content string) (int64, error) {
	uuid, err := utils.GenerateUUID()
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
