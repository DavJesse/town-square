package database

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize the DB handle
func Init(dbname string) (*sql.DB, error) {
	var err error

	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, errors.New("could not open the database")
	}

	// initialize tables
	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
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
		POSTS_TABLE_INDEX_user_id,
		POST_CATEGORIES_TABLE_CREATE,
		POST_CATEGORIES_TABLE_INDEX_post_id,
		COMMENTS_TABLE_CREATE,
		COMMENTS_TABLE_INDEX_user_id,
		COMMENTS_TABLE_INDEX_post_id,
		LIKES_TABLE_CREATE,
		LIKES_TABLE_INDEX_post_id,
		LIKES_TABLE_INDEX_comment_id,
		DISLIKES_TABLE_CREATE,
		DISLIKES_TABLE_INDEX_post_id,
		DISLIKES_TABLE_INDEX_comment_id,
	}

	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
