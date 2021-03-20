package model

type File struct {
	Model
	Size    int64  `json:"size"`
	Name    string `json:"name"`
	Mime    string `json:"mime"`
	Preview string `json:"preview"`

	DiskSpaceId uint
	DiskSpace   DiskSpace `gorm:"foreignKey:DiskSpaceId"`
}
