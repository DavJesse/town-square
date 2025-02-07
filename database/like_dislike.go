package database

import (
	"fmt"
)

// LikePost adds a like for a post and removes any existing dislike for the same post.
func LikePost(userID int, postID string) error {
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

	// Check if the user already liked the post
	var likeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)`, userID, postID).Scan(&likeExists)
	if err != nil {
		return fmt.Errorf("failed to check if like exists: %w", err)
	}

	// If already liked, just return
	if likeExists {
		return nil
	}

	// Remove any existing dislike for the same post
	_, err = tx.Exec(`DELETE FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to remove dislike: %w", err)
	}

	// Insert the like
	_, err = tx.Exec(`INSERT INTO likes (user_id, post_id) VALUES (?, ?)`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to insert post like: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

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
