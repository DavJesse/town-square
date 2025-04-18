package database

import (
	"encoding/json"
	"log"

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

func PostsFilterByUser(userID int) ([]models.PostWithUsername, error) {
	query := `
	SELECT 
    p.uuid, 
    p.title, 
    p.content, 
    p.media, 
    p.created_at, 
    u.username,
    u.first_name,
    u.last_name,
    u.image,
    COALESCE((SELECT COUNT(*) FROM likes l WHERE l.post_id = p.uuid), 0) AS likes_count,
    COALESCE((SELECT COUNT(*) FROM dislikes d WHERE d.post_id = p.uuid), 0) AS dislikes_count,
    COALESCE(
        json_group_array(
            json_object(
                'id', c.uuid, 
                'content', c.content, 
                'created_at', c.created_at, 
                'username', cu.username
            )
        ) FILTER (WHERE c.uuid IS NOT NULL), 
        '[]'
    ) AS comments
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN comments c ON c.post_id = p.uuid
	LEFT JOIN users cu ON cu.id = c.user_id
	WHERE p.user_id = ?  -- 🔽 filter by specific user ID
	GROUP BY p.uuid, u.username, u.first_name, u.last_name, u.image
	ORDER BY p.created_at DESC;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostWithUsername
	var commentsJSON string

	for rows.Next() {
		var post models.PostWithUsername
		err := rows.Scan(
			&post.UUID,
			&post.Title,
			&post.Content,
			&post.Media,
			&post.CreatedAt,
			&post.CreatorUsername,
			&post.CreatorFirstName,
			&post.CreatorLastName,
			&post.CreatorImage,
			&post.LikesCount,
			&post.DislikesCount,
			&commentsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Convert the JSON string into a slice of CommentWithCreator
		if err := json.Unmarshal([]byte(commentsJSON), &post.Comments); err != nil {
			log.Println("Error parsing comments JSON:", err)
			post.Comments = []models.CommentWithCreator{} // Ensure it's an empty slice on error
		}

		posts = append(posts, post)
	}

	return posts, nil
}
