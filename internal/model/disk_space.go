package model

const (
	DiskSpaceDefaultSpace int64 = 1e+8
)

type DiskSpace struct {
	Model

	UserOwnerId uint `gorm:"index"`

	OverallSpace  int64 // in bytes
	FreeSpace     int64 // in bytes
	OccupiedSpace int64 // in bytes
}
