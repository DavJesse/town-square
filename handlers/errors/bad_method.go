package errors

import (
	"net/http"
)

// Serves Method Not Allowed error page
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	// Populate error message and code
	ErrorResponse.Code = http.StatusMethodNotAllowed
	ErrorResponse.Issue = "Method Not Allowed"

	ErrorHandler(w, r)
}
