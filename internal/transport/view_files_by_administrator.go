package transport

import (
	"errors"
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
)

type (
	ViewFilesByAdminRequest struct {
		UserName string `json:"user_name"`
	}
	ViewFilesByAdminResponse struct {
		DiskSpace *model.DiskSpace `json:"disk_space"`
	}
)

func (req *ViewFilesByAdminRequest) Validate() error {
	if req.UserName == "" {
		return errs.NewStack(errors.New("username is empty"))
	}
	return nil
}
