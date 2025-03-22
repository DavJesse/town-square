package errors

import (
	"net/http"
)

// Serves Bad Request error page
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Populate error message and code
	ErrorResponse.Code = http.StatusBadRequest
	ErrorResponse.Issue = "Bad Request"

	ErrorHandler(w, r)
}
