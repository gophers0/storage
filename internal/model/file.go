package model

type File struct {
	Model

	Name     string
	Size     int
	FileMime string // file type

	DiskSpaceId int
	CatalogId   int // file catalog id
}
