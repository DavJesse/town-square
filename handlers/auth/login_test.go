package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler_GET(t *testing.T) {
	// Create a request with GET method
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// Check if the response body contains the expected form elements
	expectedBody := []string{
		"<form",
		"method=\"POST\"",
		"action=\"/login\"",
		"input",
		"name=\"email_username\"",
		"name=\"password\"",
		"type=\"submit\"",
	}

	for _, expected := range expectedBody {
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("handler response does not contain %s", expected)
		}
	}
}

func TestLoginHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("PUT", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestLoginHandler_FormFields(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `name="email_username"`) {
		t.Error("login form does not contain email/username field")
	}
	if !strings.Contains(body, `name="password"`) {
		t.Error("login form does not contain password field")
	}
}
