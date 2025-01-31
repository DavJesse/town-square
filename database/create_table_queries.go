package database

const (
	USERS_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username STRING UNIQUE NOT NULL,
		email STRING UNIQUE NOT NULL,
		password STRING NOT NULL,
		bio STRING,
		image STRING,
		session_id STRING
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	USERS_TABLE_INDEX_username   = `CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);`
	USERS_TABLE_INDEX_email      = `CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);`
	USERS_TABLE_INDEX_session_id = `CREATE INDEX IF NOT EXISTS idx_users_session_id ON users (session_id);`

	SESSION_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS sessions (
		session_id STRING PRIMARY KEY,
		expiry TIMESTAMP NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	SESSION_TABLE_INDEX_user_id = `CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);`
	SESSION_TABLE_INDEX_expiry  = `CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions (expiry);`

	CATEGORIES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY,
		name STRING UNIQUE NOT NULL
	);`
	CATEGORIES_TABLE_INDEX_name = `CREATE INDEX IF NOT EXISTS idx_categories_name ON categories (name);`

	POSTS_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY,
		title STRING NOT NULL,
		content STRING NOT NULL,
		media STRING,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	POSTS_TABLE_INDEX_user_id = `CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);`

	POST_CATEGORIES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`
	POST_CATEGORIES_TABLE_INDEX_post_id     = `CREATE INDEX IF NOT EXISTS idx_post_categories_post_id ON post_categories (post_id);`
	POST_CATEGORIES_TABLE_INDEX_category_id = `CREATE INDEX IF NOT EXISTS idx_post_categories_category_id ON post_categories (category_id);`

	COMMENTS_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		content STRING NOT NULL,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	COMMENTS_TABLE_INDEX_user_id = `CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);`
	COMMENTS_TABLE_INDEX_post_id = `CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);`

	LIKES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		post_id INTEGER,
		comment_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
	);`
	LIKES_TABLE_INDEX_post_id    = `CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes (post_id);`
	LIKES_TABLE_INDEX_comment_id = `CREATE INDEX IF NOT EXISTS idx_likes_comment_id ON likes (comment_id);`

	DISLIKES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS dislikes (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		post_id INTEGER,
		comment_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
	);`
	DISLIKES_TABLE_INDEX_post_id    = `CREATE INDEX IF NOT EXISTS idx_dislikes_post_id ON dislikes (post_id);`
	DISLIKES_TABLE_INDEX_comment_id = `CREATE INDEX IF NOT EXISTS idx_dislikes_comment_id ON dislikes (comment_id);`
)
