package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// Serves Not Found error page
func NotFoundHandler(w http.ResponseWriter) {
	// Set relevant header
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")

	// Set parameters of error
	hitch.Code = http.StatusNotFound
	hitch.Issue = "Not Found!"

	// Encode data to json
	if err := json.NewEncoder(w).Encode(hitch); err != nil {
		http.Error(w, `{"code":500, "issue":"internal server error"}`, http.StatusInternalServerError)
		log.Println("JSON ENCODING ERROR: ", err)
	}
}
