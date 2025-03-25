package errors

import (
	"net/http"
)

// Serves Internal Server Error page
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Populate error message and code
	ErrorResponse.Code = http.StatusInternalServerError
	ErrorResponse.Issue = "Internal Server Error"

	ErrorHandler(w, r)
}
