package errors

import (
	"net/http"

	"forum/models"
)

// Serves Bad Request error page
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Populate error message and code
	errResponse := models.WebError{
		Code:  http.StatusBadRequest,
		Issue: "Bad Request",
	}

	ErrorHandler(errResponse, w, r)
}
