package model

const (
	DiskSpaceDefaultSpace uint = 1e+8
)

type DiskSpace struct {
	Model

	UserOwnerId uint `gorm:"index"`

	OverallSpace  uint // in bytes
	FreeSpace     uint // in bytes
	OccupiedSpace uint // in bytes
}
