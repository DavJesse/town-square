package errors

import (
	"encoding/json"
	"net/http"

	"forum/models"
)

var ErrorResponse models.WebError

// Serves Not Found error page
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(ErrorResponse)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorResponse.Code = http.StatusNotFound
	ErrorResponse.Issue = "Not Found"

	ErrorHandler(w, r)
}
