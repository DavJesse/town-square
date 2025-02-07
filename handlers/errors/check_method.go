package errors

import (
	"fmt"
	"log"
	"net/http"
)

// checks the http method
func CheckHTTPMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	// Check if the method matches the allowed method
	if r.Method == allowedMethod {
		return true
	}

	// If the method does not match, log the error
	log.Printf("Invalid method: %v. Expected: %v", r.Method, allowedMethod)

	// Render error page for method not allowed
	respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method Not Allowed! Expected method: %v", allowedMethod))
	return false
}
