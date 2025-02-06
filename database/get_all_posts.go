package database

import (
	"forum/models"
)

// fetches all posts from the database with the creator's names and the number of likes and dislikes
func GetAllPosts() ([]models.PostWithUsername, error) {
	// SQL query to fetch posts along with like and dislike counts
	query := `
		SELECT p.uuid, p.title, p.content, p.media, p.created_at, u.username,
			COALESCE(COUNT(l.id), 0) AS likes_count, 
			COALESCE(COUNT(d.id), 0) AS dislikes_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.uuid
		LEFT JOIN dislikes d ON d.post_id = p.uuid
		GROUP BY p.uuid, u.username
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostWithUsername
	for rows.Next() {
		var post models.PostWithUsername
		err := rows.Scan(&post.UUID, &post.Title, &post.Content, &post.Media, &post.CreatedAt, &post.Creator, &post.LikesCount, &post.DislikesCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
