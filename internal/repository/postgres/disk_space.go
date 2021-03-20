package postgres

import (
	"errors"
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

func (r *Repo) FindOrCreateUserDiskSpace(user_id uint) (*model.DiskSpace, error) {

	mux.Lock()
	defer mux.Unlock()

	diskSpace := &model.DiskSpace{
		UserOwnerId:   user_id,
		OverallSpace:  model.DiskSpaceDefaultSpace,
		FreeSpace:     model.DiskSpaceDefaultSpace,
		OccupiedSpace: 0,
		ReservedSpace: 0,
	}

	if err := r.DB.
		Where(model.DiskSpace{UserOwnerId: user_id}).
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

func (r *Repo) ReserveDiskSpace(user_id uint, volume uint) (*model.DiskSpace, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	diskSpace := &model.DiskSpace{}
	if err = r.DB.
		Where(model.DiskSpace{UserOwnerId: user_id}).
		First(diskSpace).Error; err != nil {
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

func (r *Repo) CancelReserveDiskSpace(user_id uint, volume uint) (*model.DiskSpace, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	diskSpace := &model.DiskSpace{}
	if err = r.DB.
		Where(model.DiskSpace{UserOwnerId: user_id}).
		First(diskSpace).Error; err != nil {
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

func (r *Repo) AproveReserveDiskSpace(user_id uint, volume uint) (*model.DiskSpace, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	diskSpace := &model.DiskSpace{}
	if err = r.DB.
		Where(model.DiskSpace{UserOwnerId: user_id}).
		First(diskSpace).Error; err != nil {
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
