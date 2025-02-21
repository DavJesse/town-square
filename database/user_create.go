package database

import "forum/models"

func CreateUser(username, email, password string) (int64, error) {
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := db.Exec(query, username, email, password)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func CreateNewUser(user models.User) (int64, error) {
	query := `INSERT INTO users (username, email, password, bio, image) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.Email, user.Password, user.Bio, user.Image)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
