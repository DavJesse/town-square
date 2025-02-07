package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// module level variables
var (
	db  *sql.DB
	err error
)

// Initialize the DB handle
func Init(dbname string) error {
	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return errors.New("could not open the database")
	}
	err = enableForeignKeys(db)
	if err != nil {
		return errors.New("error enebling foreign key constraints")
	}

	err = enableWALMode(db)
	if err != nil {
		return errors.New("error enebling foreign key constraints")
	}

	// initialize tables
	err = createTables(db)
	if err != nil {
		return err
	}

	err = InitCategories()
	if err != nil {
		return err
	}

	return nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		USERS_TABLE_CREATE,
		USERS_TABLE_INDEX_username,
		USERS_TABLE_INDEX_email,
		USERS_TABLE_INDEX_session_id,
		SESSION_TABLE_CREATE,
		SESSION_TABLE_INDEX_user_id,
		SESSION_TABLE_INDEX_expiry,
		CATEGORIES_TABLE_CREATE,
		CATEGORIES_TABLE_INDEX_name,
		POSTS_TABLE_CREATE,
		POSTS_TABLE_INDEX_uuid,
		POSTS_TABLE_INDEX_user_id,
		POST_CATEGORIES_TABLE_CREATE,
		POST_CATEGORIES_TABLE_INDEX_post_id,
		POST_CATEGORIES_TABLE_INDEX_category_id,
		COMMENTS_TABLE_CREATE,
		COMMENTS_TABLE_INDEX_uuid,
		COMMENTS_TABLE_INDEX_user_id,
		COMMENTS_TABLE_INDEX_post_id,
		LIKES_TABLE_CREATE,
		LIKES_TABLE_INDEX_user_id,
		LIKES_TABLE_INDEX_post_id,
		LIKES_TABLE_INDEX_comment_id,
		DISLIKES_TABLE_CREATE,
		DISLIKES_TABLE_INDEX_user_id,
		DISLIKES_TABLE_INDEX_post_id,
		DISLIKES_TABLE_INDEX_comment_id,
	}

	for _, query := range tables {
		_, err = db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() {
	db.Close()
}

// SQLite does not enforce foreign key constraints by default
func enableForeignKeys(db *sql.DB) error {
	// enable foreign key constraints for this connection
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign key constraints: %w", err)
	}
	return nil
}

/*
* WAL (Write-Ahead Logging) Mode
* In WAL mode, SQLite allows for greater concurrency by allowing multiple readers while there is a single writer.
* This mode helps in reducing the chances of database locks when there are multiple concurrent read and write operations
 */
func enableWALMode(db *sql.DB) error {
	_, err := db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
	}
	return nil
}
