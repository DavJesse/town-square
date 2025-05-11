package database

import (
	"encoding/json"
	"html"
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
                'uuid', c.uuid, 
                'content', c.content,
				'post_id', c.post_id,
				'creator_first_name', cu.first_name,
				'creator_last_name', cu.last_name,
				'creator_username', cu.username,
				'creator_image', cu.image, 
                'created_at', strftime('%Y-%m-%dT%H:%M:%SZ', c.created_at), 
                'username', cu.username
            )
        ) FILTER (WHERE c.uuid IS NOT NULL), 
        '[]'
    ) AS comments
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN comments c ON c.post_id = p.uuid
	LEFT JOIN users cu ON cu.id = c.user_id
	WHERE p.user_id = ?  -- filter by specific user ID
	GROUP BY p.uuid, u.username, u.first_name, u.last_name, u.image
	ORDER BY p.created_at DESC;
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostWithUsername
	var commentsJSON string
	var postContent string

	for rows.Next() {
		var post models.PostWithUsername
		err := rows.Scan(
			&post.UUID,
			&post.Title,
			&postContent,
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

		// Unescape HTML entities
		post.Content = html.UnescapeString(postContent)
		commentsJSON = html.UnescapeString(commentsJSON)

		// Convert the JSON string into a slice of CommentWithCreator
		if err := json.Unmarshal([]byte(commentsJSON), &post.Comments); err != nil {
			log.Println("Error parsing comments JSON:", err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
