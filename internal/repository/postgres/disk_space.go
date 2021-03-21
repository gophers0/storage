package postgres

import (
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
)

func (r *Repo) FindDiskSpace(id uint) (*model.DiskSpace, error) {
	mux.RLock()
	defer mux.RUnlock()

	diskSpace := &model.DiskSpace{}
	if err := r.DB.
		Where("id = ?", id).
		First(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}

func (r *Repo) FindOrCreateUserDiskSpace(userId uint) (*model.DiskSpace, error) {

	mux.Lock()
	defer mux.Unlock()

	diskSpace := &model.DiskSpace{
		UserOwnerId:   userId,
		OverallSpace:  model.DiskSpaceDefaultSpace,
		FreeSpace:     model.DiskSpaceDefaultSpace,
		OccupiedSpace: 0,
	}

	if err := r.DB.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Preload("Files").
		Where(model.DiskSpace{UserOwnerId: userId}).
		FirstOrCreate(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return diskSpace, nil
}

func (r *Repo) DeleteDiskSpace(id uint) (*model.DiskSpace, []*model.File, error) {
	mux.Lock()
	defer mux.Unlock()

	var err error
	diskSpace := &model.DiskSpace{}

	// create a transaction
	tx := r.DB.Begin()

	err = tx.
		Where("id = ?", id).
		Delete(diskSpace).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, errs.NewStack(err)
	}

	files := []*model.File{}
	err = tx.
		Where(model.File{DiskSpaceId: diskSpace.ID}).
		Delete(files).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, errs.NewStack(err)
	}

	tx.Commit()

	return diskSpace, files, err
}

func (r *Repo) FillDiskSpace(userId uint, volume int64) (*model.DiskSpace, error) {
	mux.Lock()
	defer mux.Unlock()
	diskSpace := &model.DiskSpace{}
	if err := r.DB.
		Where(model.DiskSpace{UserOwnerId: userId}).
		First(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	if volume > diskSpace.FreeSpace {
		return nil, errs.NotAvailableFreeSpace
	}

	diskSpace.FreeSpace -= volume
	diskSpace.OccupiedSpace += volume

	if err := r.DB.Save(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return diskSpace, nil
}
