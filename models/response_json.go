package models

type Response struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Redirect string `json:"redirect"`
	Data     any    `json:"data"`
}
