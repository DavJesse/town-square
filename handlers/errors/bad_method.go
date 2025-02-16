package errors

import (
	"html/template"
	"log"
	"net/http"

	utils "forum/utils"
)

// Serves Bad Request error page
func MethodNotAllowedHandler(w http.ResponseWriter) {
	// Construct absolute path to error.html
	tmplPath, err := utils.GetTemplatePath("error.html")
	if err != nil {
		http.Error(w, "error page unavailable", http.StatusNotFound)
		log.Printf("TEMPLATE AVAILABILITY ERROR: %v", err)
		return
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("TEMPLATE PARSING ERROR:", err)
		http.Error(w, "Could not load template, error page unavailable", http.StatusInternalServerError)
		return
	}

	// Set relevant headers
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set parameters of error
	hitch.Code = http.StatusBadRequest
	hitch.Issue = "Method Not Allowed!"

	// Execute bad request template, handle emerging errors
	err = tmpl.Execute(w, hitch)
	if err != nil {
		http.Error(w, "Could not execute error template", http.StatusInternalServerError)
		log.Println("Template EXECUTION ERROR: ", err)
	}
}
