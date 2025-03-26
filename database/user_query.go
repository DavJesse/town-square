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
	query := `SELECT id, first_name, last_name, age, gender, username, email, bio, image, created_at 
	          FROM users WHERE session_id = ?`

	var user models.User
	var bio, image sql.NullString // Use sql.NullString for nullable fields
	var firstName, lastName, gender sql.NullString
	var age sql.NullInt64 // Use sql.NullInt64 for nullable integers

	err := db.QueryRow(query, UUID).Scan(
		&user.ID,
		&firstName,
		&lastName,
		&age,
		&gender,
		&user.Username,
		&user.Email,
		&bio,
		&image,
		&user.CreatedAt,
	)

	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return models.User{}, err
	}

	// Convert NullString to string (use empty string if NULL)
	if firstName.Valid {
		user.FirstName = firstName.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}
	if gender.Valid {
		user.Gender = gender.String
	}
	if age.Valid {
		user.Age = int(age.Int64)
	}
	if bio.Valid {
		user.Bio = bio.String
	}
	if image.Valid {
		user.Image = image.String
	}

	return user, nil
}

func GetUserbyID(userID int) (models.User, error) {
	query := `SELECT id, username, email, bio, image, created_at FROM users WHERE id = ?`

	var user models.User
	var bio, image sql.NullString // Use sql.NullString for nullable fields

	err := db.QueryRow(query, userID).Scan(
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
