package errors

import (
	"html/template"
	"log"
	"net/http"

	"forum/models"
	"forum/utils"
)

var hitch models.WebError

// Centralized function for rendering errors
func respondWithError(w http.ResponseWriter, statusCode int, issue string) {
	// Set up the error message
	hitch := models.WebError{
		Code:  statusCode,
		Issue: issue,
	}

	// Get template path
	tmplPath, err := utils.GetTemplatePath("error.html")
	if err != nil {
		log.Println("Error finding template file:", err)
		http.Error(w, "Error finding template file", http.StatusInternalServerError)
		return
	}

	// Parse and render the template
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Error parsing template file:", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Write the error response
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, hitch); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
