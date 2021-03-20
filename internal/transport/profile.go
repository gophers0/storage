package transport

import "github.com/gophers0/storage/internal/model"

type (
	ProfileRequest struct {
	}
	ProfileResponse struct {
		DiskSpace   *model.DiskSpace `json:"disk_space"`
		Files       []*model.File    `json:"files"`
		TrashFiles  []*model.File    `json:"trash_files"`
		SharedFiles []*model.File    `json:"shared_files"`
	}
)

func (req *ProfileRequest) Validate() error {
	return nil
}
