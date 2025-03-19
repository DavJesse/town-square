package errors

import (
	"net/http"

	"forum/models"
)

// Serves Internal Server Error page
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Populate error message and code
	errResponse := models.WebError{
		Code:  http.StatusInternalServerError,
		Issue: "Internal Server Error",
	}

	ErrorHandler(errResponse, w, r)
}
