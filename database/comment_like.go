package database

import (
	"fmt"
)

// LikeComment adds a like for a comment and removes any existing dislike for the same comment.
func LikeComment(userID int, commentID string) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if the user already liked the comment
	var likeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND comment_id = ?)`, userID, commentID).Scan(&likeExists)
	if err != nil {
		return fmt.Errorf("failed to check if like exists: %w", err)
	}

	// If already liked, just return
	if likeExists {
		return nil
	}

	// Remove any existing dislike for the same comment
	_, err = tx.Exec(`DELETE FROM dislikes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to remove dislike: %w", err)
	}

	// Insert the like
	_, err = tx.Exec(`INSERT INTO likes (user_id, comment_id) VALUES (?, ?)`, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to insert comment like: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
