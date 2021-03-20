package model

const (
	AccessRightIdOwner int = 1

	AccessRightIdRead   int = 11 // see in catalog, open in app
	AccessRightIdLoad   int = 12
	AccessRightIdDelete int = 21

	AccessRightIdGiveRead        int = 31
	AccessRightIdGiveDelete      int = 32
	AccessRightIdUploadToCatalog int = 1
)

type UserAccessRight struct {
	Model

	UserId            int
	AccessRightTypeId int // AccessRight constant
}
