package database

import (
	"fmt"

	"forum/models"
)

// Fetches all posts from the database where the user ID matches the provided ID
func GetPostsByUserID(id int) ([]models.Post, error) {
	query := `
		SELECT 
			p.uuid, 
			p.title, 
			p.content, 
			p.media, 
			p.created_at
		FROM posts p
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.UUID, &post.Title, &post.Content, &post.Media, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		posts = append(posts, post)
	}

	// check for any row iteration error
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return posts, nil
}
