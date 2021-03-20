package model

const (
	AccessRightIdOwner uint = 1

	AccessRightIdSee    uint = 11 // see in catalog
	AccessRightIdRead   uint = 12 // see in open in app
	AccessRightIdLoad   uint = 13
	AccessRightIdDelete uint = 21

	AccessRightIdGiveRead        uint = 31
	AccessRightIdGiveDelete      uint = 32
	AccessRightIdUploadToCatalog uint = 51
)

type UserAccessRight struct {
	Model

	FileId            uint
	UserId            uint `gorm:"index"` // we do not store user data
	AccessRightTypeId uint // AccessRight constant

	File File `gorm:"foreignKey:FileId"`
}
