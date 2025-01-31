CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username STRING UNIQUE NOT NULL,
    email STRING UNIQUE NOT NULL,
    password STRING NOT NULL,
    bio STRING,
    image STRING,
    session_id STRING
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_session_id ON users (session_id);

CREATE TABLE sessions (
    session_id STRING PRIMARY KEY,
    expiry TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_expiry ON sessions (expiry);

CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name STRING UNIQUE NOT NULL
);
CREATE INDEX idx_categories_name ON categories (name);

CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    title STRING NOT NULL,
    content STRING NOT NULL,
    media STRING,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_posts_user_id ON posts (user_id);

CREATE TABLE post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
CREATE INDEX idx_post_categories_post_id ON post_categories (post_id);
CREATE INDEX idx_post_categories_category_id ON post_categories (category_id);

CREATE TABLE comments (
    id INTEGER PRIMARY KEY,
    content STRING NOT NULL,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_comments_user_id ON comments (user_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);

CREATE TABLE likes (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
);
CREATE INDEX idx_likes_post_id ON likes (post_id);
CREATE INDEX idx_likes_comment_id ON likes (comment_id);

CREATE TABLE dislikes (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    CHECK (post_id IS NOT NULL OR comment_id IS NOT NULL)
);
CREATE INDEX idx_dislikes_post_id ON dislikes (post_id);
CREATE INDEX idx_dislikes_comment_id ON dislikes (comment_id);