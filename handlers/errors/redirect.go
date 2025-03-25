package errors

import (
	"encoding/json"
	"net/http"

	"forum/models"
)

// Custom handler for redirects via JSON
func RedirectHandler(w http.ResponseWriter, redirectPath string) {
	jsonResponse := models.Response{
		Code:     http.StatusSeeOther,
		Message:  "Redirecting",
		Redirect: redirectPath,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusSeeOther)
	json.NewEncoder(w).Encode(jsonResponse)
}
