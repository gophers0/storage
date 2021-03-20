package postgres

import (
	"errors"
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
	"path/filepath"
)

func (r *Repo) FindFile(id uint) (*model.File, error) {
	mux.RLock()
	defer mux.RUnlock()

	file := &model.File{}
	if err := r.DB.
		Where("id = ?", id).
		First(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return file, nil
}

func (r *Repo) FindFiles(ids []uint) ([]*model.File, error) {
	mux.RLock()
	defer mux.RUnlock()

	files := []*model.File{}
	if err := r.DB.
		Where("id in (?)", ids).
		First(files).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return files, nil
}

func (r *Repo) FindDiskFiles(disk_space_id uint) ([]*model.File, error) {
	mux.RLock()
	defer mux.RUnlock()

	catalogs := []*model.File{}
	if err := r.DB.
		Where(model.File{DiskSpaceId: disk_space_id}).
		Find(catalogs).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalogs, nil
}

func (r *Repo) CreateFile(name string, size, disk_space_id, catalog_id uint) (*model.File, error) {
	var err error

	// check dick and catalog
	_, err = r.FindDiskSpace(disk_space_id)
	if err != nil {
		return nil, errs.NewStack(errors.New("Disk space does not exists!"))
	}

	_, err = r.FindCatalog(catalog_id)
	if err != nil {
		return nil, errs.NewStack(errors.New("Catalog does not exists!"))
	}

	mux.Lock()
	defer mux.Unlock()

	file := &model.File{}
	file.Name = name
	file.FileMime = filepath.Ext(name)
	file.Size = size
	file.DiskSpaceId = disk_space_id
	file.CatalogId = catalog_id

	if err := r.DB.
		Create(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return file, nil
}

func (r *Repo) DeleteFile(id uint) (*model.File, error) {
	var err error
	_, err = r.FindFile(id)
	if err != nil {
		return nil, errs.NewStack(errors.New("File does not exists!"))
	}

	mux.Lock()
	defer mux.Unlock()

	file := &model.File{}
	if err := r.DB.
		Where("id = ?", id).
		First(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return file, nil
}
