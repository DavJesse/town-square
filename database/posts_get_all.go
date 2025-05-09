package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"forum/models"
	"forum/utils"
)

// fetches all posts from the database with the creator's names and the number of likes and dislikes
func GetAllPosts() ([]models.PostWithUsername, error) {
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
	    COALESCE(json_group_array(json_object(
									'uuid', c.uuid,
									'content', c.content,
									'post_id', c.post_id,
									'creator_first_name', cu.first_name,
									'creator_last_name', cu.last_name,
									'creator_username', cu.username,
									'creator_image', cu.image,
									'created_at', strftime('%Y-%m-%dT%H:%M:%SZ', c.created_at),
									'likes_count', COALESCE((SELECT COUNT(*) FROM likes l WHERE l.comment_id = c.uuid), 0),
    								'dislikes_count', COALESCE((SELECT COUNT(*) FROM dislikes d WHERE d.comment_id = c.uuid), 0)
									))
	             FILTER (WHERE c.uuid IS NOT NULL), '[]') AS comments
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN comments c ON c.post_id = p.uuid
	LEFT JOIN users cu ON cu.id = c.user_id
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
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetLikedPostsByUser(userID int) ([]models.PostWithCategories, error) {
	query := `
		SELECT 
			p.uuid,
			u.first_name,
			u.last_name, 
			u.username,
			u.image, 
			p.title, 
			p.content, 
			p.media, 
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
			&post.CreatorFirstName,
			&post.CreatorLastName,
			&post.CreatorUsername,
			&post.CreatorImage,
			&post.Title,
			&post.Content,
			&post.Media,
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

		// Fetch comments for this post
		comments, err := GetCommentsForPost(post.UUID)
		if err != nil {
			log.Printf("Failed to fetch comments for post %s: %v", post.UUID, err)
			comments = []models.CommentWithCreator{}
		}
		post.Comments = comments

		likedPosts = append(likedPosts, post)
	}

	return likedPosts, nil
}

// Fetch liked posts from a particular category
func FetchLikedPostsPerCategory(categoryID int, userID int) ([]models.PostWithUsername, error) {
	// Check if the category exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", categoryID).Scan(&exists)
	if err != nil {
		log.Println("Error checking if category exists:", err)
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("category with ID %d does not exist", categoryID)
	}

	query := `
		SELECT  
		    p.uuid, 
		    u.first_name, 
		    u.last_name, 
		    u.username, 
		    u.image, 
		    p.title, 
		    p.content, 
		    p.media, 
		    p.created_at, 
		    COALESCE(l.likes_count, 0) AS likes_count, 
		    COALESCE(d.dislikes_count, 0) AS dislikes_count, 
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
		INNER JOIN post_categories pc ON p.uuid = pc.post_id
		INNER JOIN users u ON u.id = p.user_id
		INNER JOIN likes ul ON ul.post_id = p.uuid AND ul.user_id = ?
		LEFT JOIN (
		    SELECT post_id, COUNT(*) AS likes_count  
		    FROM likes  
		    GROUP BY post_id
		) l ON p.uuid = l.post_id
		LEFT JOIN (
		    SELECT post_id, COUNT(*) AS dislikes_count  
		    FROM dislikes  
		    GROUP BY post_id
		) d ON p.uuid = d.post_id
		LEFT JOIN comments c ON c.post_id = p.uuid
		LEFT JOIN users cu ON cu.id = c.user_id
		WHERE pc.category_id = ?
		GROUP BY p.uuid
		ORDER BY p.created_at DESC;
	`

	rows, err := db.Query(query, userID, categoryID)
	if err != nil {
		log.Println("Error querying posts by category ID:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostWithUsername
	for rows.Next() {
		var post models.PostWithUsername
		var commentsJSON string
		var media sql.NullString

		err := rows.Scan(
			&post.UUID,
			&post.CreatorFirstName,
			&post.CreatorLastName,
			&post.CreatorUsername,
			&post.CreatorImage,
			&post.Title,
			&post.Content,
			&media,
			&post.CreatedAt,
			&post.LikesCount,
			&post.DislikesCount,
			&commentsJSON,
		)
		if err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}
		post.Media = ""
		if media.Valid {
			post.Media = media.String
		}
		if err := json.Unmarshal([]byte(commentsJSON), &post.Comments); err != nil {
			log.Printf("DATABASE ERROR: %v", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return posts, nil
}

func GetCommentsForPost(postUUID string) ([]models.CommentWithCreator, error) {
	query := `
	SELECT 
		c.uuid, 
		c.content,
		c.post_id,  
		c.created_at,
		u.first_name,
		u.last_name,
		u.username,
		u.image,
		COALESCE((SELECT COUNT(*) FROM likes l WHERE l.comment_id = c.uuid), 0) AS likes_count,
		COALESCE((SELECT COUNT(*) FROM dislikes d WHERE d.comment_id = c.uuid), 0) AS dislikes_count
	FROM comments c
	INNER JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.created_at ASC
`

	rows, err := db.Query(query, postUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %v", err)
	}
	defer rows.Close()

	var comments []models.CommentWithCreator

	for rows.Next() {
		var comment models.CommentWithCreator
		err := rows.Scan(
			&comment.UUID,
			&comment.Content,
			&comment.PostID,
			&comment.CreatedAt,
			&comment.CreatorFirstName,
			&comment.CreatorLastName,
			&comment.CreatorUsername,
			&comment.CreatorImage,
			&comment.LikesCount,
			&comment.DislikesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}

		// Convert time to EAT
		eatTime, err := utils.ConvertToEAT(comment.CreatedAt.String())
		if err == nil {
			comment.CreatedAt = eatTime
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
