package model

type File struct {
	Model
	Size uint
	Name string
	Mime string // file type

	DiskSpaceId uint
	DiskSpace   DiskSpace `gorm:"foreignKey:DiskSpaceId"`
}
