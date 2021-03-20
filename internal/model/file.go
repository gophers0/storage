package model

type File struct {
	Model
	Size     uint
	Name     string
	FileMime string // file type

	DiskSpaceId uint
	CatalogId   uint // file catalog id

	DiskSpace DiskSpace `gorm:"foreignKey:DiskSpaceId"`
	Catalog   Catalog   `gorm:"foreignKey:CatalogId"`
}
