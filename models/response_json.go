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

type PostResponse struct {
	Categories []Category `json:"categories"`
	IsLogged   bool       `json:"is_logged"`
	ProfPic    string     `json:"prof_pic"`
	Message    string     `json:"message"`
	Code       int        `json:"code,omitempty"`
}

type PostEngagementResponse struct {
	LikesCount   int `json:"like_count"`
	DislikeCount int `json:"dislike_count"`
}

type DislikeResponse struct {
	Success       bool   `json:"success"`
	DislikesCount int    `json:"dislikesCount"`
	Message       string `json:"message"`
}
