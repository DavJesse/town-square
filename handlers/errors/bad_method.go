package errors

import (
	"net/http"

	"forum/models"
)

// Serves Method Not Allowed error page
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	// Populate error message and code
	errResponse := models.WebError{
		Code:  http.StatusMethodNotAllowed,
		Issue: "MethodNotAllowed",
	}

	ErrorHandler(errResponse, w, r)
}
