package models

type Response struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Redirect string `json:"redirect"`
	Data     any    `json:"data"`
}

type LoginResponse struct {
	EmailUsername string `json:"email_username"`
	Password      string `json:"password"`
}
