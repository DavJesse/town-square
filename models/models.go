package models

type WebError struct {
	Code  int
	Issue string
}

type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Username   string
	Password   string
	Bio        string
	ProfilePic string
}
