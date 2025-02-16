package auth

// import (
// 	"html/template"
// 	"net/http"

// 	"forum/handlers/errors"
// )

// func ForgotPassword(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFiles("web/templates/forgot_password.html")
// 	if err != nil {
// 		http.Error(w, "Failed to load Success template", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := tmpl.Execute(w, r); err != nil {
// 		errors.InternalServerErrorHandler(w)
// 		return
// 	}
// }
