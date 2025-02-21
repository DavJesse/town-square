package database

import (
	"database/sql"
	"forum/models" // Replace with your actual path
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func TestSearchPosts(t *testing.T) {
	// 1. Setup: Use an in-memory SQLite database for testing.
	db, err := sql.Open("sqlite3", ":memory:") // In-memory database
	if err != nil {
		t.Fatal("Error opening test database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS posts (
                        uuid VARCHAR(255) PRIMARY KEY,
                        title VARCHAR(255),
                        content TEXT,
            media VARCHAR(255),
            user_id INT,
            created_at TIMESTAMP
                );
        `)
	if err != nil {
		t.Fatal("Error creating test table:", err)
	}

	//Insert sample data
	_, err = db.Exec(`
                INSERT INTO posts (uuid, title, content, media, user_id, created_at) VALUES ('1', 'Test Post 1', 'This is a test post.', 'test.jpg', 1, datetime('now')), ('2', 'Another Test', 'This is another test post.', '', 2, datetime('now')), ('3', 'No Match Here', 'This post should not be returned.', '', 3, datetime('now'));
        `)
	if err != nil {
		t.Fatal("Error inserting test data:", err)
	}

	// 2. Test Cases (same as before)
	testCases := []struct {
		name        string
		query       string
		expectedLen int
		expected    []models.Post // Or []SearchResult if you created that struct
	}{
		// ... (Test cases remain the same)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			posts, err := SearchPosts(tc.query)
			if err != nil {
				t.Errorf("SearchPosts returned an error: %v", err)
			}

			if len(posts) != tc.expectedLen {
				t.Errorf("Expected %d results, got %d", tc.expectedLen, len(posts))
			}

			//Important:  Check the *contents* of the returned posts.
			for i, post := range posts {
				if i < len(tc.expected) { //Check bounds to avoid panics
					if post.Title != tc.expected[i].Title || post.Content != tc.expected[i].Content {
						t.Errorf("Result %d: Expected Title '%s' and Content '%s', got '%s' and '%s'", i, tc.expected[i].Title, tc.expected[i].Content, post.Title, post.Content)
					}
				}
			}
		})
	}

	//3. Teardown: Clean up the test database (if necessary).  The `defer db.Close()` handles closing the connection.
	_, err = db.Exec("DROP TABLE posts;") //Clean up after the test
	if err != nil {
		log.Println("Error dropping table:", err) //Log the error, but don't fail the test
	}
}
