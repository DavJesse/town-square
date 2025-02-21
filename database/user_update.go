package database

func UpdateUser(userID int, bio, image string) error {
	query := `UPDATE users SET bio = ?, image = ? WHERE id = ?`
	_, err := db.Exec(query, bio, image, userID)
	return err
}
