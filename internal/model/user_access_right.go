package model

const (
	AccessRightIdOwner uint = 1

	AccessRightIdRead   uint = 11 // see in catalog, open in app
	AccessRightIdLoad   uint = 12
	AccessRightIdDelete uint = 21

	AccessRightIdGiveRead        uint = 31
	AccessRightIdGiveDelete      uint = 32
	AccessRightIdUploadToCatalog uint = 1
)

type UserAccessRight struct {
	Model

	UserId            uint
	AccessRightTypeId uint // AccessRight constant
}
