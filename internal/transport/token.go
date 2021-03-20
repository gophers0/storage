package transport

import "github.com/gophers0/users/internal/model"

type CheckTokenRequest struct {
	Token  string `json:"token"`
	UserId int    `json:"user_id"`
}

func (ctr *CheckTokenRequest) Validate() error {
	return nil
}

type CheckTokenResponse struct {
	Success bool           `json:"success"`
	Session *model.Session `json:"session"`
}
