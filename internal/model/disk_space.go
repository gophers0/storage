package model

type DiskSpace struct {
	Model

	UserOwnerId int

	OverallVolume  uint // in bytes
	FreeVolume     uint // in bytes
	OccupiedVolume uint // in bytes
	ReservedVolume uint // in bytes
}
