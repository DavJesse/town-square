package errors

import (
	"encoding/json"
	"net/http"

	"forum/models"
)

// Serves Not Found error page
func ErrorHandler(errResponse models.WebError, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(errResponse)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	errResponse := models.WebError{
		Code:  http.StatusNotFound,
		Issue: "Not Found",
	}

	ErrorHandler(errResponse, w, r)
}
