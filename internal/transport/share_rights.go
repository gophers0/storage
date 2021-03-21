package transport

import (
	"errors"
	"fmt"

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
	if req.UserName == "" {
		return errs.NewStack(errors.New("username is empty"))
	}
	if req.FileId <= 0 {
		return errs.NewStack(errors.New(fmt.Sprintf("invalid file id: %d", req.FileId)))
	}
	return nil
}
