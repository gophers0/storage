package transport

import (
	"errors"
	"github.com/gophers0/storage/pkg/errs"
)

type (
	ShareReadRightRequest struct {
		UserName string `json:"user_name"`
		FileId   int    `json:"file_id"`
	}
	ShareReadRightResponse struct {
	}
)

func (req *ShareReadRightRequest) Validate() error {
	var err error
	if len(req.UserName) < 3 {
		err = errs.NewStack(errors.New("The username is incorrect."))
	}
	if req.FileId < 0 {
		err = errs.NewStack(errors.New("The file id is incorrect."))
	}
	return err
}
