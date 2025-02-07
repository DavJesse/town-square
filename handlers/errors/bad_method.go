package errors

import (
	"html/template"
	"log"
	"net/http"

	"forum/utils"
)

// Serves Bad Request error page
func BadRequestHandler(w http.ResponseWriter) {
	// Construct absolute path to error.html
	tmplPath, err := utils.GetTemplatePath("error.html")
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
