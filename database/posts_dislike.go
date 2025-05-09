package database

import "fmt"

// DislikePost adds a dislike for a post and removes any existing like for the same post.
func DislikePost(userID int, postID string) error {
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

	// Check if the user already disliked the post
	var dislikeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM dislikes WHERE user_id = ? AND post_id = ?)`, userID, postID).Scan(&dislikeExists)
	if err != nil {
		return fmt.Errorf("failed to check if dislike exists: %w", err)
	}

	if dislikeExists {
		// If the post is already disliked, undislike it (remove the dislike)
		_, err = tx.Exec(`DELETE FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			return fmt.Errorf("failed to remove dislike: %w", err)
		}
	} else {
		// Remove any existing like before adding the dislike
		_, err = tx.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			return fmt.Errorf("failed to remove like: %w", err)
		}

		// Insert the dislike
		_, err = tx.Exec(`INSERT INTO dislikes (user_id, post_id) VALUES (?, ?)`, userID, postID)
		if err != nil {
			return fmt.Errorf("failed to insert post dislike: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetPostDislikesCount(postID string) (int, error) {
	var dislikesCount int

	// Retrieve the count of dislikes for the post
	err := db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE post_id = ?`, postID).Scan(&dislikesCount)
	if err != nil {
		return 0, err
	}
	return dislikesCount, nil
}
