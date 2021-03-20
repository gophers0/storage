package postgres

import (
	"errors"
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
	"sync"
)

var mux = sync.Mutex{} // simple slow way to solve race condition problem in ReservedDiskSpace

func (r *Repo) FindDiskSpace(id uint) (*model.DiskSpace, error) {
	diskSpace := &model.DiskSpace{}
	if err := r.DB.
		Where("id = ?", id).
		First(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}

func (r *Repo) FindUserDiskSpace(user_id uint) (*model.DiskSpace, error) {
	diskSpace := &model.DiskSpace{}
	if err := r.DB.
		Where(model.DiskSpace{UserOwnerId: user_id}).
		First(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}

func (r *Repo) CreateDiskSpace(user_id uint) (*model.DiskSpace, error) {
	diskSpace := &model.DiskSpace{}

	diskSpace.UserOwnerId = user_id
	diskSpace.OverallSpace = model.DiskSpaceDefaultSpace
	diskSpace.FreeSpace = diskSpace.OverallSpace
	diskSpace.OccupiedSpace = 0
	diskSpace.ReservedSpace = 0

	if err := r.DB.
		Where(model.DiskSpace{UserOwnerId: user_id}).
		FirstOrCreate(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}

func (r *Repo) DeleteDiskSpace(user_id uint) (*model.DiskSpace, error) {
	var err error
	diskSpace := &model.DiskSpace{}

	// create a transaction
	tx := r.DB.Begin()

	err = tx.
		Where(model.DiskSpace{UserOwnerId: user_id}).
		Delete(diskSpace).Error
	if err != nil {
		tx.Rollback()
		return diskSpace, errs.NewStack(err)
	}

	catalogs := []*model.Catalog{}
	err = tx.
		Where(model.Catalog{DiskSpaceId: diskSpace.ID}).
		Delete(catalogs).Error
	if err != nil {
		tx.Rollback()
		return diskSpace, errs.NewStack(err)
	}

	files := []*model.File{}
	err = tx.
		Where(model.File{DiskSpaceId: diskSpace.ID}).
		Delete(files).Error
	if err != nil {
		tx.Rollback()
		return diskSpace, errs.NewStack(err)
	}

	tx.Commit()

	return diskSpace, err
}

func (r *Repo) ReservedDiskSpace(user_id uint, volume uint) (*model.DiskSpace, error) {
	mux.Lock()
	defer mux.Unlock()

	var err error
	diskSpace, err := r.FindUserDiskSpace(user_id)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	if volume > diskSpace.FreeSpace {
		return nil, errs.NewStack(errors.New("Not enough free space"))
	}
	diskSpace.FreeSpace -= volume
	diskSpace.ReservedSpace += volume

	if err := r.DB.
		Update(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}

func (r *Repo) CancelReservedDiskSpace(user_id uint, volume uint) (*model.DiskSpace, error) {
	var err error
	diskSpace, err := r.FindUserDiskSpace(user_id)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	if volume > diskSpace.ReservedSpace {
		return nil, errs.NewStack(errors.New("Less memory reserved than expected"))
	}
	diskSpace.FreeSpace += volume
	diskSpace.ReservedSpace -= volume

	if err := r.DB.
		Update(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}

func (r *Repo) AproveReservedDiskSpace(user_id uint, volume uint) (*model.DiskSpace, error) {
	var err error
	diskSpace, err := r.FindUserDiskSpace(user_id)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	if volume > diskSpace.ReservedSpace {
		return nil, errs.NewStack(errors.New("Less memory reserved than expected"))
	} else if volume > diskSpace.FreeSpace {
		return nil, errs.NewStack(errors.New("Less memory free than expected"))
	}
	diskSpace.ReservedSpace -= volume
	diskSpace.OccupiedSpace += volume

	if err := r.DB.
		Update(diskSpace).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return diskSpace, nil
}
