package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"forum/models"
)

// initializes the categories table with the predefined categories
func InitCategories() error {
	categories := []string{
		"miscellaneous",
		"agriculture",
		"arts",
		"education",
		"lifestyle",
		"technology",
		"culture",
		"science",
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, category := range categories {
		query := `INSERT OR IGNORE INTO categories (name) VALUES (?)`
		_, err := tx.Exec(query, category)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", category, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// fetchCategories retrieves all categories from the database
func FetchCategories() ([]models.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// validateCategories checks if the provided category IDs exist in the database
func ValidateCategories(categoryIDs []int) error {
	query := `SELECT COUNT(*) FROM categories WHERE id IN (?` + strings.Repeat(",?", len(categoryIDs)-1) + `)`
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return err
	}

	if count != len(categoryIDs) {
		return fmt.Errorf("one or more categories do not exist")
	}

	return nil
}

// FetchCategoryPostsWithID retrieves all posts associated with a given category ID
func FetchCategoryPostsWithID(categoryID int) ([]models.PostWithUsername, error) {
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

	// Fetch posts for the category
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

	rows, err := db.Query(query, categoryID)
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
