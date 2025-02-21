package database

import (
	"fmt"
	"strings"

	"forum/models"
	"forum/utils"
)

// fetches all posts from the database with the creator's names and the number of likes and dislikes
func GetAllPosts() ([]models.PostWithUsername, error) {
	query := `
		SELECT p.uuid, p.title, p.content, p.media, p.created_at, u.username,
		    COALESCE((SELECT COUNT(*) FROM likes l WHERE l.post_id = p.uuid), 0) AS likes_count,
		    COALESCE((SELECT COUNT(*) FROM dislikes d WHERE d.post_id = p.uuid), 0) AS dislikes_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
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

func GetLikedPostsByUser(userID string) ([]models.PostWithCategories, error) {
	query := `
		SELECT 
			p.uuid, 
			p.title, 
			p.content, 
			p.media, 
			u.username, 
			p.user_id, 
			p.created_at,
			GROUP_CONCAT(DISTINCT c.name) AS category_names,
			COUNT(DISTINCT l.id) AS likes_count, 
			COUNT(DISTINCT dl.id) AS dislikes_count
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_categories pc ON p.uuid = pc.post_id
		LEFT JOIN categories c ON pc.category_id = c.id
		INNER JOIN likes l ON p.uuid = l.post_id  -- Only fetch posts the user liked
		LEFT JOIN dislikes dl ON p.uuid = dl.post_id
		WHERE l.user_id = ?
		GROUP BY p.uuid, p.title, p.content, p.media, u.username, p.user_id, p.created_at
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch liked posts: %v", err)
	}
	defer rows.Close()

	var likedPosts []models.PostWithCategories

	for rows.Next() {
		var post models.PostWithCategories
		var categoryNames string

		err := rows.Scan(
			&post.UUID,
			&post.Title,
			&post.Content,
			&post.Media,
			&post.Username,
			&post.UserID,
			&post.CreatedAt,
			&categoryNames,
			&post.LikesCount,
			&post.DislikesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning liked posts: %v", err)
		}

		// Convert created_at time to East African Time
		eatTime, err := utils.ConvertToEAT(post.CreatedAt.String())
		if err == nil {
			post.CreatedAt = eatTime
		}

		// Convert category names to a slice
		if categoryNames != "" {
			post.Categories = strings.Split(categoryNames, ",")
		} else {
			post.Categories = []string{}
		}

		likedPosts = append(likedPosts, post)
	}

	return likedPosts, nil
}
