package postgres

import (
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
)

func (r *Repo) FindUserAccessRight(id uint) (*model.UserAccessRight, error) {
	mux.RLock()
	defer mux.RUnlock()

	userAccessRight := &model.UserAccessRight{}
	if err := r.DB.
		Where("id = ?", id).
		First(userAccessRight).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return userAccessRight, nil
}

func (r *Repo) FindUserAccessRights(userId uint) ([]*model.UserAccessRight, error) {
	mux.RLock()
	defer mux.RUnlock()

	var userAccessRights []*model.UserAccessRight
	if err := r.DB.
		Where(model.UserAccessRight{UserId: userId}).
		Find(&userAccessRights).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return userAccessRights, nil
}

func (r *Repo) CreateUserAccessRight(userId, fileId, accessRightId uint) (*model.UserAccessRight, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	file := &model.File{}
	if err = r.DB.
		Where("id = ?", fileId).
		First(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	userAccessRight := &model.UserAccessRight{
		UserId:            userId,
		FileId:            fileId,
		AccessRightTypeId: accessRightId,
	}

	if err := r.DB.Debug().Where(userAccessRight).FirstOrCreate(userAccessRight).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return userAccessRight, nil
}

func (r *Repo) DeleteUserAccessRight(userId, fileId, accessRightId uint) (*model.UserAccessRight, error) {

	mux.Lock()
	defer mux.Unlock()

	userAccessRight := &model.UserAccessRight{}

	if err := r.DB.
		Where(&model.UserAccessRight{
			UserId:            userId,
			FileId:            fileId,
			AccessRightTypeId: accessRightId,
		}).
		Delete(userAccessRight).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return userAccessRight, nil
}
