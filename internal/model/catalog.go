package model

type Catalog struct {
	Model

	DiskSpaceId uint

	ParentCatalogId uint // zero for root directories
	Name            string
	DiskSpace       DiskSpace `gorm:"foreignKey:DiskSpaceId"`
}
