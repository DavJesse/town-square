package database

const (
	USERS_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		first_name STRING NOT NULL,
		last_name STRING NOT NULL,
        age INTEGER CHECK (age >= 13 AND age <= 130),
        gender STRING CHECK (gender IN ('male', 'female', 'other')),
		username STRING UNIQUE NOT NULL,
        email STRING UNIQUE NOT NULL,
		password STRING NOT NULL,  
		bio STRING,
		image STRING,
		session_id STRING,
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
		uuid STRING PRIMARY KEY,
		title STRING NOT NULL,
		content STRING NOT NULL,
		media STRING,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	POSTS_TABLE_INDEX_uuid    = `CREATE INDEX IF NOT EXISTS idx_posts_uuid ON posts (uuid);`
	POSTS_TABLE_INDEX_user_id = `CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);`

	POST_CATEGORIES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS post_categories (
		post_id STRING NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(uuid),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`
	POST_CATEGORIES_TABLE_INDEX_post_id     = `CREATE INDEX IF NOT EXISTS idx_post_categories_post_id ON post_categories (post_id);`
	POST_CATEGORIES_TABLE_INDEX_category_id = `CREATE INDEX IF NOT EXISTS idx_post_categories_category_id ON post_categories (category_id);`

	COMMENTS_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS comments (
		uuid STRING PRIMARY KEY,
		content STRING NOT NULL,
		post_id STRING NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(uuid),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	COMMENTS_TABLE_INDEX_uuid    = `CREATE INDEX IF NOT EXISTS idx_comments_uuid ON comments (uuid);`
	COMMENTS_TABLE_INDEX_user_id = `CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);`
	COMMENTS_TABLE_INDEX_post_id = `CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);`

	LIKES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		post_id STRING,  -- Nullable post_id
		comment_id STRING,  -- Nullable comment_id
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(uuid),
		FOREIGN KEY (comment_id) REFERENCES comments(uuid),
		CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL),  -- At least one must be non-null
		UNIQUE (user_id, post_id),  -- Ensure the user can only like the same post once
		UNIQUE (user_id, comment_id) -- Ensure the user can only like the same comment once
	);`
	LIKES_TABLE_INDEX_user_id    = `CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes (user_id);`
	LIKES_TABLE_INDEX_post_id    = `CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes (post_id);`
	LIKES_TABLE_INDEX_comment_id = `CREATE INDEX IF NOT EXISTS idx_likes_comment_id ON likes (comment_id);`

	DISLIKES_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS dislikes (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		post_id STRING,  -- Nullable post_id
		comment_id STRING,  -- Nullable comment_id
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(uuid),
		FOREIGN KEY (comment_id) REFERENCES comments(uuid),
		CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL),  -- At least one must be non-null
		UNIQUE (user_id, post_id),  -- Ensure the user can only dislike the same post once
		UNIQUE (user_id, comment_id) -- Ensure the user can only dislike the same comment once
	);`
	DISLIKES_TABLE_INDEX_user_id    = `CREATE INDEX IF NOT EXISTS idx_dislikes_user_id ON dislikes (user_id);`
	DISLIKES_TABLE_INDEX_post_id    = `CREATE INDEX IF NOT EXISTS idx_dislikes_post_id ON dislikes (post_id);`
	DISLIKES_TABLE_INDEX_comment_id = `CREATE INDEX IF NOT EXISTS idx_dislikes_comment_id ON dislikes (comment_id);`

	MESSAGES_TABLE_CREATE = `
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sender_id INTEGER NOT NULL,
            receiver_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            timestamp DATETIME NOT NULL,
            FOREIGN KEY (sender_id) REFERENCES users(id),
            FOREIGN KEY (receiver_id) REFERENCES users(id)
        )
    `
	USER_STATUS_TABLE_CREATE = `CREATE TABLE IF NOT EXISTS user_status (
		user_id INTEGER PRIMARY KEY,
		is_online BOOLEAN DEFAULT false,
		last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status_message TEXT DEFAULT NULL,
		device_info TEXT DEFAULT NULL,
		notification_preferences TEXT DEFAULT 'all',
		typing_status BOOLEAN DEFAULT false,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_user_status_last_activity ON user_status(last_activity);
    `
)
