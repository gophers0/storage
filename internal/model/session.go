package model

type Session struct {
	Model
	User   User   `json:"user"`
	UserID uint   `json:"-"`
	Token  string `json:"token"`
}
