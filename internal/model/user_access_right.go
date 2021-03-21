package model

const (
	AccessRightIdOwner uint = 1

	AccessRightIdRead   uint = 11 // see in open in app
	AccessRightIdDelete uint = 21

	AccessRightIdGiveRead        uint = 31
	AccessRightIdGiveDelete      uint = 32
	AccessRightIdUploadToCatalog uint = 51
)

type UserAccessRight struct {
	Model
	FileId            uint `json:"file_id"`
	UserId            uint `json:"user_id" gorm:"index"` // we do not store user data
	AccessRightTypeId uint `json:"access_right_type_id"` // AccessRight constant
}
