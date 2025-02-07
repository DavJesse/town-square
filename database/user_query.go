package database

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/utils"
)

func GetUserByEmailOrUsername(email, username string) (models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = ? OR username = ?`
	var user models.User
	err := db.QueryRow(query, email, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func VerifyUser(email, password string) (bool, error) {
	// User email or username to get users' full credentials
	user, err := GetUserByEmailOrUsername(email, email)
	if err != nil {
		return false, fmt.Errorf("user does not exist: %v", err)
	}

	// Compare provided password with the stored hashed password
	return utils.MatchPasswords(user.Password, password)
}

// GetUserbySessionID function
func GetUserbySessionID(UUID string) (models.User, error) {
	// fmt.Println("Session ID:", UUID)
	query := `SELECT id, username, email, bio, image, created_at FROM users WHERE session_id = ?`

	var user models.User
	var bio, image sql.NullString // Use sql.NullString for nullable fields

	err := db.QueryRow(query, UUID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&bio,   // Scan into NullString
		&image, // Scan into NullString
		&user.CreatedAt,
	)

	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return models.User{}, err
	}

	// Convert NullString to string, using empty string if NULL
	if bio.Valid {
		user.Bio = bio.String
	}
	if image.Valid {
		user.Image = image.String
	}

	return user, nil
}
