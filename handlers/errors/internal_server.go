package errors

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/models"
)

var hitch models.WebError

// Serves Internal Server Error page
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Set relevant headers
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "./web/templates/index.html")

	// Set parameters of error
	hitch.Code = http.StatusInternalServerError
	hitch.Issue = "Internal Server Error!"

	// Encode data to json
	if err := json.NewEncoder(w).Encode(hitch); err != nil {
		http.Error(w, `{"code":500, "issue":"internal server error"}`, http.StatusInternalServerError)
		log.Println("JSON ENCODING ERROR: ", err)
	}
}
