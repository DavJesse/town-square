package database

import "fmt"

// DislikePost adds a dislike for a post and removes any existing like for the same post.
func DislikePost(userID int, postID string) error {
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

	// Check if the user already disliked the post
	var dislikeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM dislikes WHERE user_id = ? AND post_id = ?)`, userID, postID).Scan(&dislikeExists)
	if err != nil {
		return fmt.Errorf("failed to check if dislike exists: %w", err)
	}

	// If already disliked, just return
	if dislikeExists {
		return nil
	}

	// Remove any existing like for the same post
	_, err = tx.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}

	// Insert the dislike
	_, err = tx.Exec(`INSERT INTO dislikes (user_id, post_id) VALUES (?, ?)`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to insert post dislike: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
