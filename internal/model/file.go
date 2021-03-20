package model

type File struct {
	Model
	Size     uint
	Name     string
	FileMime string // file type

	DiskSpaceId uint
	CatalogId   uint // file catalog id
}
