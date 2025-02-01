package database

import (
	"database/sql"

	"forum/models"
)

func GetCommentsByPostUUID(db *sql.DB, postID int) ([]models.Comment, error) {
	query := `SELECT uuid, content, user_id, created_at FROM comments WHERE post_id = ?`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.UUID, &comment.Content, &comment.UserID, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
