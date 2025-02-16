package database

import (
	"fmt"
	"strings"

	"database/sql"
	"forum/models"
	"log"
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
func FetchCategoryPostsWithID(categoryID int) ([]models.Post, error) {
	query := `
        SELECT p.uuid, p.title, p.content, p.media, p.user_id, p.created_at
        FROM posts p
        JOIN post_categories pc ON p.uuid = pc.post_id
        WHERE pc.category_id = ?`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		log.Println("Error querying posts by category ID:", err)
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
