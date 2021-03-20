package transport

import (
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (req *LoginRequest) Validate() error {
	if req.Login == "" {
		return errs.RequestValidationError.AddInfo("login is empty")
	}
	if req.Password == "" {
		return errs.RequestValidationError.AddInfo("password is empty")
	}
	return nil
}

type LoginResponse struct {
	Code  int         `json:"code"`
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}
