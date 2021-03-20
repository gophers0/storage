package model

const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

type User struct {
	Model
	Login    string `json:"login"`
	Password string `json:"-"`
	Role     string `json:"role"`
}
