package database

import (
	"forum/models"
)

func GetPostByUUID(postID int) (models.Post, error) {
	query := `SELECT uuid, title, content, media, user_id, created_at FROM posts WHERE id = ?`
	var post models.Post
	err := db.QueryRow(query, postID).Scan(&post.UUID, &post.Title, &post.Content, &post.Media, &post.UserID, &post.CreatedAt)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}
