package database

import "fmt"

// LikePost adds a like for a post and removes any existing dislike for the same post.
func LikePost(userID int, postID string) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if the user already liked the post
	var likeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)`, userID, postID).Scan(&likeExists)
	if err != nil {
		return fmt.Errorf("failed to check if like exists: %w", err)
	}

	if likeExists {
		// If the post is already liked, unlike it (remove the like)
		_, err = tx.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			return fmt.Errorf("failed to remove like: %w", err)
		}
	} else {
		// Remove any existing dislike before adding the like
		_, err = tx.Exec(`DELETE FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			return fmt.Errorf("failed to remove dislike: %w", err)
		}

		// Insert the like
		_, err = tx.Exec(`INSERT INTO likes (user_id, post_id) VALUES (?, ?)`, userID, postID)
		if err != nil {
			return fmt.Errorf("failed to insert post like: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
