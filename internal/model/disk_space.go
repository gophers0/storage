package model

const (
	DiskSpaceDefaultSpace int = 1e+8
)

type DiskSpace struct {
	Model

	UserOwnerId uint

	OverallSpace  int // in bytes
	FreeSpace     int // in bytes
	OccupiedSpace int // in bytes
	ReservedSpace int // in bytes
}
