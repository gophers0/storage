package model

type Catalog struct {
	Model

	DiskSpaceId uint

	IsMainCatalog   bool // is it main catalog in disk space?
	ParentCatalogId uint // zero for main catalog

	DiskSpace DiskSpace `gorm:"foreignKey:DiskSpaceId"`
}
