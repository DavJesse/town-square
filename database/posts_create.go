package database

import (
	"fmt"
	"strings"
)

// Create inserts data into the specified table using provided field names and values.
func CreatePostWithCategories(userID int, title, content, media string, categoryIDs []int) (int64, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	postUUID, err := GenerateUUID()
	if err != nil {
		return 0, fmt.Errorf("failed to generate a random uuid: %w", err)
	}

	postValues := []interface{}{postUUID, userID, title, content, media}

	// prepare the SQL query for the post
	postPlaceholders := make([]string, len(postValues))
	for i := range postValues {
		postPlaceholders[i] = "?"
	}

	postQuery := fmt.Sprintf("INSERT INTO posts (uuid, user_id, title, content, media) VALUES (%s)",
		strings.Join(postPlaceholders, ", "),
	)

	// execute the post query
	postResult, err := tx.Exec(postQuery, postValues...)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	// get the last inserted post ID
	postID, err := postResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve post ID: %w", err)
	}

	// insert post categories
	for _, categoryID := range categoryIDs {
		categoryValues := []interface{}{postUUID, categoryID}

		// prepare the SQL query for the category
		categoryPlaceholders := make([]string, len(categoryValues))
		for i := range categoryValues {
			categoryPlaceholders[i] = "?"
		}

		categoryQuery := fmt.Sprintf("INSERT INTO post_categories (post_id, category_id) VALUES (%s)",
			strings.Join(categoryPlaceholders, ", "),
		)

		// execute the category query
		_, err := tx.Exec(categoryQuery, categoryValues...)
		if err != nil {
			return 0, fmt.Errorf("failed to associate category with post: %w", err)
		}
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return postID, nil
}
