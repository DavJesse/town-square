package utils

import (
	"forum/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Passwordhash(user *models.User) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password", err)
		return
	}
	user.Password = string(hashedPassword)
}
