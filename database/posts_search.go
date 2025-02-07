package database

import (
	"database/sql"
	"forum/models"
	"log"
)

// SearchPosts searches for posts matching the query(title or content)
func SearchPosts(query string) ([]models.Post, error) {
	query = "%" + query + "%"
	sqlQuery := `
        SELECT uuid, title, content, media, user_id, created_at
        FROM posts
        WHERE title LIKE ? OR content LIKE ?`
	rows, err := db.Query(sqlQuery, query, query)
	if err != nil {
		log.Println("Error querying posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var media sql.NullString
		err := rows.Scan(&post.UUID, &post.Title, &post.Content, &media, &post.UserID, &post.CreatedAt)
		if err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}
		if media.Valid {
			post.Media = media.String
		} else {
			post.Media = ""
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return posts, nil
}
