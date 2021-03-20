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
		Find(files).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return files, nil
}

func (r *Repo) FindDiskFiles(diskSpaceId uint) ([]*model.File, error) {
	mux.RLock()
	defer mux.RUnlock()

	catalogs := []*model.File{}
	if err := r.DB.
		Where(model.File{DiskSpaceId: diskSpaceId}).
		Find(catalogs).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalogs, nil
}

func (r *Repo) CreateFile(name string, size, diskSpaceId uint) (*model.File, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	// check dick
	diskSpace := &model.DiskSpace{}
	if err = r.DB.
		Where("id = ?", diskSpaceId).
		First(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	if err != nil {
		return nil, errs.NewStack(errors.New("Disk space does not exists!"))
	}

	file := &model.File{
		Name:        name,
		Size:        size,
		DiskSpaceId: diskSpaceId,
	}
	file.FileMime = filepath.Ext(name)

	if err := r.DB.
		Create(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return file, nil
}

func (r *Repo) DeleteFile(id uint) (*model.File, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	file := &model.File{}
	if err = r.DB.
		Where("id = ?", id).
		First(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return file, nil
}
