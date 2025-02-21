package database

// import (
// 	"fmt"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"forum/models"
// )

// func TestIsLoggedIn(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	result := IsLoggedIn(req)

// 	if result {
// 		t.Errorf("IsLoggedIn() = %v, want %v", result, false)
// 	}
// }

// func TestIsLoggedIn_EmptySessionCookie(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Set an empty session cookie
// 	cookie := &http.Cookie{
// 		Name:  "session_id",
// 		Value: "",
// 	}
// 	req.AddCookie(cookie)

// 	result := IsLoggedIn(req)
// 	if result {
// 		t.Errorf("IsLoggedIn() returned true for empty session cookie, expected false")
// 	}
// }

// func TestIsLoggedIn_SessionNotFound(t *testing.T) {
// 	// Create a mock request with a session cookie
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	cookie := &http.Cookie{
// 		Name:  "session_id",
// 		Value: "invalid_session_id",
// 	}
// 	req.AddCookie(cookie)

// 	// Mock the GetSession function to return an error
// 	GetSession := func(sessionID string) (models.Session, error) {
// 		return models.Session{}, fmt.Errorf("session not found")
// 	}
// 	oldGetSession := GetSession
// 	defer func() { GetSession = oldGetSession }()

// 	// Call the IsLoggedIn function
// 	result := IsLoggedIn(req)

// 	// Check if the result is false
// 	if result {
// 		t.Errorf("Expected IsLoggedIn to return false when session is not found, but got true")
// 	}
// }

// func TestIsLoggedIn_ExpiredSession(t *testing.T) {
// 	// Create a mock request with a session cookie
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	cookie := &http.Cookie{
// 		Name:  "session_id",
// 		Value: "test_session_id",
// 	}
// 	req.AddCookie(cookie)

// 	// Mock the GetSession function to return an expired session
// 	GetSession := func(sessionID string) (models.Session, error) {
// 		return models.Session{
// 			SessionID: "test_session_id",
// 			UserID:    1,
// 			Expiry:    time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
// 		}, nil
// 	}
// 	oldGetSession := GetSession
// 	defer func() { GetSession = oldGetSession }()

// 	// Call the function
// 	result := IsLoggedIn(req)

// 	// Check the result
// 	if result {
// 		t.Errorf("IsLoggedIn() = %v, want %v", result, false)
// 	}
// }

// func TestIsLoggedIn_DatabaseError(t *testing.T) {
// 	// Mock the request with a valid session cookie
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	cookie := &http.Cookie{
// 		Name:  "session_id",
// 		Value: "test_session_id",
// 	}
// 	req.AddCookie(cookie)

// 	// Call the function
// 	result := IsLoggedIn(req)

// 	// Check the result
// 	if result {
// 		t.Error("Expected IsLoggedIn to return false on database error, but got true")
// 	}
// }

// func TestIsLoggedIn_NonExpiredCookie(t *testing.T) {
// 	// Create a mock request with a valid session cookie
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	cookie := &http.Cookie{
// 		Name:  "session_id",
// 		Value: "3cb056f2-3458-4137-863b-32751e68928e",
// 	}
// 	req.AddCookie(cookie)

// 	// Call the IsLoggedIn function
// 	result := IsLoggedIn(req)

// 	// Assert that the result is true
// 	if !result {
// 		t.Errorf("Expected IsLoggedIn to return true, but got false")
// 	}
// }
