package database

func DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, userID)
	return err
}
