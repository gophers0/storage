package transport

import "github.com/gophers0/storage/pkg/errs"

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (req *CreateUserRequest) Validate() error {
	if req.Login == "" {
		return errs.RequestValidationError.AddInfo("empty login")
	}
	return nil
}
