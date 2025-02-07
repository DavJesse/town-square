package utils

import (
	"forum/models"
	"net/http"
)

func ParseUserForm(r *http.Request) (*models.User, *models.WebError) {
	// Extract form data
	user := models.User{}
	user.Email = r.FormValue("email")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Password = r.FormValue("confirm_password")
	user.Bio = r.FormValue("bio")
	err := r.ParseForm()
	if err != nil {
		return nil, &models.WebError{Code: http.StatusBadRequest, Issue: "Invalid form submission"}
	}
	return &user, nil
}
