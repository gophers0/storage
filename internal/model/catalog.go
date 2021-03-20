package model

type Catalog struct {
	Model

	DiskSpaceId int

	IsMainCatalog   bool // is it main catalog in disk space?
	ParentCatalogId int  // zero for main catalog
}
