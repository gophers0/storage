package model

type File struct {
	Model
	Size int64
	Name string
	Mime string // file type

	DiskSpaceId uint
	DiskSpace   DiskSpace `gorm:"foreignKey:DiskSpaceId"`
}
