package errors

import (
	"encoding/json"
	"net/http"
)

// JSONErrorResponse sends a JSON error response with the given message and status code.
func JSONErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})
}
