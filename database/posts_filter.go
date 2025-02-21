package database

import (
	"forum/models"
)

func PostsFilterByCategory(categoryID int) ([]models.Post, error) {
	query := `
        SELECT p.uuid, p.title, p.content, p.media, p.user_id, p.created_at 
        FROM posts p
        JOIN post_categories pc ON p.uuid = pc.post_id
        WHERE pc.category_id = ?`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.UUID, &post.Title, &post.Content, &post.Media, &post.UserID, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func PostsFilterByUser(userID int) ([]models.Post, error) {
	query := `SELECT uuid, title, content, media, created_at FROM posts WHERE user_id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.UUID, &post.Title, &post.Content, &post.Media, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
