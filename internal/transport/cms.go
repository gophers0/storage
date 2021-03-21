package transport

import "github.com/gophers0/storage/internal/model"

type ListUsersRequest struct {
	Login  string `json:"login"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (req *ListUsersRequest) Validate() error {
	return nil
}

type ListUsersResponse struct {
	Count   int        `json:"count"`
	Records []UserInfo `json:"records"`
}

type UserInfo struct {
	Id        int              `json:"id"`
	CreatedAt string           `json:"created_at"`
	Login     string           `json:"login"`
	Role      string           `json:"role"`
	DiskSpace *model.DiskSpace `json:"disk_space"`
}
