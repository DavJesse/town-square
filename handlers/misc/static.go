package misc

import (
	"log"
	"net/http"
	"os"
	"strings"

	"forum/handlers/errors"
)

func Static(w http.ResponseWriter, r *http.Request) {
	// Restric requests to "GET"
	if r.Method != "GET" {
		errors.MethodNotAllowedHandler(w)
		log.Printf("METHOD ERROR: method not allowed")
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
		errors.NotFoundHandler(w)
		log.Println("DIRECTORY AVAILABILITY ERROR: directory not found")
		return
	}

	// Set predetermined path to catch malicious file tranversal
	filePath := "./web/" + path
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		errors.NotFoundHandler(w)
		log.Printf("FILE AVAILABILITY ERROR: %v", err)
		return
	}

	// Restrict access to directories
	if fileInfo.IsDir() {
		errors.NotFoundHandler(w)
		log.Println("DIRECTORY AVAILABILITY ERROR: directory not found")
		return
	}

	// Serve the requested file
	http.ServeFile(w, r, filePath)
}
