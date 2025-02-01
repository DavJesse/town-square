package database

func CreateUser(username, email, password string) (int64, error) {
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := db.Exec(query, username, email, password)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
