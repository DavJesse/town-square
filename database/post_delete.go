package database

import (
	"fmt"
)

func DeletePost(uuid string) error {
	tx, err := db.Begin() // Start a transaction
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Deleting dependent records first
	queries := []struct {
		query string
		args  []interface{}
	}{
		{"DELETE FROM comments WHERE post_id = ?", []interface{}{uuid}},
		{"DELETE FROM likes WHERE post_id = ?", []interface{}{uuid}},
		{"DELETE FROM dislikes WHERE post_id = ?", []interface{}{uuid}},
		{"DELETE FROM post_categories WHERE post_id = ?", []interface{}{uuid}},
		{"DELETE FROM posts WHERE uuid = ?", []interface{}{uuid}}, // Finally, delete the post itself
	}

	for _, q := range queries {
		if _, err := tx.Exec(q.query, q.args...); err != nil {
			tx.Rollback() // Rollback if any query fails
			return fmt.Errorf("failed to delete post dependencies: %w", err)
		}
	}

	// Commit transaction if all deletions succeed
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
