package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// Serves Not Found error page
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Get the path from the query string
	path := r.URL.Query().Get("path")

	// Set parameters of error
	hitch.Code = http.StatusNotFound
	hitch.Issue = "Not Found!"
	hitch.Path = path

	// Check if it's an API request (e.g., Fetch API)
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" || r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(hitch); err != nil {
			http.Error(w, `{"code":500, "issue":"internal server error"}`, http.StatusInternalServerError)
			log.Println("JSON ENCODING ERROR: ", err)
		}
		return
	}

	// If it's a direct browser visit, serve the SPA
	http.ServeFile(w, r, "./web/templates/index.html")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/error?path="+url.QueryEscape(r.URL.Path), http.StatusSeeOther)
}
