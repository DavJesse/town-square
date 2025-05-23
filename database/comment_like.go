package database

import (
	"fmt"
	"forum/models"
)

// LikeComment adds a like for a comment and removes any existing dislike for the same comment.
func LikeComment(userID int, commentID string) error {
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

	// Check if the user already liked the comment
	var likeExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND comment_id = ?)`, userID, commentID).Scan(&likeExists)
	if err != nil {
		return fmt.Errorf("failed to check if like exists: %w", err)
	}

	if likeExists {
		// If the comment is already liked, unlike it (remove the like)
		_, err = tx.Exec(`DELETE FROM likes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
		if err != nil {
			return fmt.Errorf("failed to remove like: %w", err)
		}
	} else {
		// Remove any existing dislike before adding the like
		_, err = tx.Exec(`DELETE FROM dislikes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
		if err != nil {
			return fmt.Errorf("failed to remove dislike: %w", err)
		}

		// Insert the like
		_, err = tx.Exec(`INSERT INTO likes (user_id, comment_id) VALUES (?, ?)`, userID, commentID)
		if err != nil {
			return fmt.Errorf("failed to insert comment like: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetCommentReactionsCounts(commentID string) (models.EngagementCount, error) {
	query := `
	SELECT
	(SELECT COUNT(*) FROM likes WHERE comment_id = ?) AS like_count,
	(SELECT COUNT(*) FROM dislikes WHERE comment_id = ?) AS dislike_count;
    `

	var engagement models.EngagementCount
	err := db.QueryRow(query, commentID, commentID).Scan(&engagement.LikesCount, &engagement.DislikesCount)
	if err != nil {
		return models.EngagementCount{}, err
	}

	return engagement, nil
}
