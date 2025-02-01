package database

import "database/sql"

func CreateComment(db *sql.DB, userID, postID int, content string) (int64, error) {
	query := `INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)`
	result, err := db.Exec(query, userID, postID, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
