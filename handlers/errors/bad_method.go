package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// Serves Bad Request error page
func MethodNotAllowedHandler(w http.ResponseWriter) {
	// Set response headers
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	// Set parameters of error
	hitch.Code = http.StatusBadRequest
	hitch.Issue = "Method Not Allowed!"

	if err := json.NewEncoder(w).Encode(hitch); err != nil {
		http.Error(w, `{"code::500, "issue":"internal server error"}`, http.StatusInternalServerError)
		log.Println("JSON ENCODING ERROR: ", err)
	}
}
