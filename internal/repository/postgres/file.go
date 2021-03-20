package postgres

import (
	"errors"
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

func (r *Repo) CreateFile(name string, size, disk_space_id uint) (*model.File, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	// check dick  !!!!!!!!!!!!!!!!!! МОЖЕТ БЫТЬ, ЧТО БЛАГОДАРЯ КЛЮЧУ МОЖНО СТЕРЕТЬ ПРОВЕРКУ
	diskSpace := &model.DiskSpace{}
	if err = r.DB.
		Where("id = ?", disk_space_id).
		First(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	if err != nil {
		return nil, errs.NewStack(errors.New("Disk space does not exists!"))
	}

	file := &model.File{
		Name:        name,
		Size:        size,
		DiskSpaceId: disk_space_id,
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
