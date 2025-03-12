package database

import (
	"database/sql"
	"log"
)

func UpdateUsersTable(db *sql.DB) error {
	// Disable foreign key constraints temporarily.
	_, err := db.Exec(`PRAGMA foreign_keys = OFF;`)
	if err != nil {
		log.Fatal("Failed to disable foreign keys:", err)
	}

	_, err = db.Exec(`
	BEGIN TRANSACTION;

	CREATE TABLE users_new (
		id INTEGER PRIMARY KEY,
		first_name STRING DEFAULT '',
		last_name STRING DEFAULT '',
		age INTEGER CHECK (age >= 13 AND age <= 130) DEFAULT 18,
		gender STRING CHECK (gender IN ('male', 'female', 'other')) DEFAULT 'other',
		username STRING UNIQUE NOT NULL,
		email STRING UNIQUE NOT NULL,
		password STRING NOT NULL,  
		bio STRING DEFAULT NULL,
		image STRING DEFAULT NULL,
		session_id STRING DEFAULT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	INSERT INTO users_new (id, first_name, last_name, age, gender, username, email, password, bio, image, session_id, created_at)
	SELECT id, '' AS first_name, '' AS last_name, 18 AS age, 'other' AS gender, username, email, password, bio, image, session_id, created_at FROM users;

	DROP TABLE users;
	ALTER TABLE users_new RENAME TO users;

	COMMIT;
`)
	if err != nil {
		return err
	}
	return nil
}
