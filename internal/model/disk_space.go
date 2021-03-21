package model

const (
	DiskSpaceDefaultSpace int64 = 1024 * 1024 * 100
)

type DiskSpace struct {
	Model
	UserOwnerId   uint   `json:"user_id" gorm:"index"`
	OverallSpace  int64  `json:"overall_space"`  // in bytes
	FreeSpace     int64  `json:"free_space"`     // in bytes
	OccupiedSpace int64  `json:"occupied_space"` // in bytes
	Files         []File `json:"files"`
}
