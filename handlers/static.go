package handlers

import (
	"net/http"
	"os"
	"strings"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	// Restric requests to "GET"
	if r.Method != "GET" {
		BadRequestHandler(w)
		return
	}

	path := r.URL.Path

	// Ensure only allowed static directories are served
	allowedPrefixes := []string{"/static/", "/css/", "/js/", "/images/"}
	valid := false
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(path, prefix) {
			valid = true
			break
		}
	}

	if !valid {
		NotFoundHandler(w)
		return
	}

	// Set predetermined path to catch malicious file tranversal
	filePath := "./web/" + path
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		NotFoundHandler(w)
		return
	}

	// Restrict access to directories
	if fileInfo.IsDir() {
		NotFoundHandler(w)
		return
	}

	// Serve the requested file
	http.ServeFile(w, r, filePath)
}
