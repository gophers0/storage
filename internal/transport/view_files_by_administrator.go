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
		Files     []*model.File    `json:"files"`
	}
)

func (req *ViewFilesByAdminRequest) Validate() error {
	var err error
	if len(req.UserName) < 3 {
		err = errs.NewStack(errors.New("The username is incorrect."))
	}
	return err
}
