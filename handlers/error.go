package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"forum/models"
)

var hitch models.WebError

// Serves Bad Request error page
func BadRequestHandler(w http.ResponseWriter) {
	// Construct absolute path to error.html
	tmplPath, err := GetTemplatePath("error.html")
	if err != nil {
		log.Printf("Could not find template file: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Template parsing failed:", err)
		http.Error(w, "Could not load template, error page unavailable", http.StatusInternalServerError)
		return
	}

	// Set relevant headers
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set parameters of error
	hitch.Code = http.StatusBadRequest
	hitch.Issue = "Bad Request!"

	// Execute bad request template, handle emerging errors
	err = tmpl.Execute(w, hitch)
	if err != nil {
		http.Error(w, "Could not execute error template, error page unavailable", http.StatusInternalServerError)
		log.Println("Error executing template: ", err)
	}
}

// Serves Internal Server Error page
func InternalServerErrorHandler(w http.ResponseWriter) {
	// Construct absolute path to error.html
	tmplPath, err := GetTemplatePath("error.html")
	if err != nil {
		log.Printf("Could not find template file: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Could not load template, error page unavailable", http.StatusInternalServerError)
		return
	}

	// Set relevant headers
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set parameters of error
	hitch.Code = http.StatusInternalServerError
	hitch.Issue = "Internal Server Error!"

	// Execute internal server error template, handle emerging errors
	err = tmpl.Execute(w, hitch)
	if err != nil {
		http.Error(w, "Could not execute error template, error page unavailable", http.StatusInternalServerError)
		log.Println("Error executing template: ", err)
	}
}

// Serves Not Found error page
func NotFoundHandler(w http.ResponseWriter) {
	// Construct absolute path to error.html
	tmplPath, err := GetTemplatePath("error.html")
	if err != nil {
		http.Error(w, "Could not find template file", http.StatusInternalServerError)
		log.Println("Could not find template file: ", err)
		return
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Could not load template, error page unavailable", http.StatusInternalServerError)
		return
	}

	// Set relevant header
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set parameters of error
	hitch.Code = http.StatusNotFound
	hitch.Issue = "Not Found!"

	// Execute not found error template, handle emerging errors
	err = tmpl.Execute(w, hitch)
	if err != nil {
		http.Error(w, "Could not execute error template, error page unavailable", http.StatusInternalServerError)
		log.Println("Error executing template: ", err)
	}
}

func GetTemplatePath(templateFile string) (string, error) {
	// catch empty template files
	if templateFile == "" {
		return "", fmt.Errorf("template file name cannot be empty")
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up until we find the project root
	dir := wd
	for {
		// Construct path to template, check if constructed path exists
		templatePath := filepath.Join(dir, "web", "templates", templateFile)
		if _, err := os.Stat(templatePath); err == nil {
			return templatePath, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir { // Stop if we reach the root
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("template file not found: %s", templateFile)
}

//method checker
func MethodCheck(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	// Check if the method matches the allowed method
	if r.Method == allowedMethod {
		return true
	}

	// If the method does not match, log the error
	log.Printf("Invalid method: %v. Expected: %v", r.Method, allowedMethod)

	// Render error page for method not allowed
	respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method Not Allowed! Expected method: %v", allowedMethod))
	return false
}

// Centralized function for rendering errors
func respondWithError(w http.ResponseWriter, statusCode int, issue string) {
	// Set up the error message
	hitch := models.WebError{
		Code:  statusCode,
		Issue: issue,
	}

	// Get template path
	tmplPath, err := GetTemplatePath("error.html")
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
