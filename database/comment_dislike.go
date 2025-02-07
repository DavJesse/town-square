package database

import "fmt"

// DislikeComment adds a dislike for a comment and removes any existing like for the same comment.
func DislikeComment(userID int, commentID string) error {
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

	// Check if the user already disliked the comment
	var dislikeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM dislikes WHERE user_id = ? AND comment_id = ?)`, userID, commentID).Scan(&dislikeExists)
	if err != nil {
		return fmt.Errorf("failed to check if dislike exists: %w", err)
	}

	// If already disliked, just return
	if dislikeExists {
		return nil
	}

	// Remove any existing like for the same comment
	_, err = tx.Exec(`DELETE FROM likes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}

	// Insert the dislike
	_, err = tx.Exec(`INSERT INTO dislikes (user_id, comment_id) VALUES (?, ?)`, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to insert comment dislike: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
