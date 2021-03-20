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
		BaseResponse
	}
)

func (req *ShareReadRightRequest) Validate() error {
	var err error
	if req.UserName == "" {
		err = errs.NewStack(errors.New("username is empty"))
	}
	if req.FileId <= 0 {
		err = errs.NewStack(errors.New("invalid file id"))
	}
	return err
}
