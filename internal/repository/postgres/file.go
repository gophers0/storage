package postgres

import (
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
	"path/filepath"
)

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

func (r *Repo) FindFile(name string, dSpaceId uint) (*model.File, error) {
	file := &model.File{}

	if err := r.DB.Model(file).Where("disk_space_id = ? AND name = ?", dSpaceId, name).First(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (r *Repo) FindDiskFiles(diskSpaceId uint) ([]*model.File, error) {
	mux.RLock()
	defer mux.RUnlock()

	catalogs := []*model.File{}
	if err := r.DB.
		Where(model.File{DiskSpaceId: diskSpaceId}).
		Find(&catalogs).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalogs, nil
}

func (r *Repo) CreateFile(name, mime string, size, diskSpaceId uint) (*model.File, error) {
	mux.Lock()
	defer mux.Unlock()

	file := &model.File{
		Name:        name,
		Size:        size,
		Mime:        mime,
		DiskSpaceId: diskSpaceId,
	}
	file.Mime = filepath.Ext(name)

	if err := r.DB.Create(file).Error; err != nil {
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
