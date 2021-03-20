package transport

import "github.com/gophers0/storage/pkg/errs"

type UploadFileReq struct {
	Name string `form:"name" json:"name"`
	Size int    `form:"size" json:"size"`
	Mime string `form:"mime" json:"mime"`
}

func (req *UploadFileReq) Validate() error {
	if req.Name == "" {
		return errs.RequestValidationError.AddInfo("filename is empty")
	}
	if req.Size == 0 {
		return errs.RequestValidationError.AddInfo("file is empty")
	}
	if req.Mime == "" {
		return errs.RequestValidationError.AddInfo("mime type is empty")
	}

	return nil
}

type UploadFileResponse struct {
	BaseResponse
	FreeSpaceAvailable int64  `json:"free_space"`
	OccupiedSpace      int64  `json:"occupied_space"`
	FileName           string `json:"file_name"`
	FileMime           string `json:"file_mime"`
	FileID             uint   `json:"file_id"`
}
