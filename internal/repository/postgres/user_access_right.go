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

func (r *Repo) FindUserAccessRights(user_id uint) ([]*model.UserAccessRight, error) {
	mux.RLock()
	defer mux.RUnlock()

	userAccessRights := []*model.UserAccessRight{}
	if err := r.DB.
		Where(model.UserAccessRight{UserId: user_id}).
		First(userAccessRights).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return userAccessRights, nil
}

func (r *Repo) CreateUserAccessRight(user_id, file_id, access_right_id uint) (*model.UserAccessRight, error) {
	var err error

	mux.Lock()
	defer mux.Unlock()

	file := &model.File{}
	if err = r.DB.
		Where("id = ?", file_id).
		First(file).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	userAccessRight := &model.UserAccessRight{
		UserId:            user_id,
		FileId:            file_id,
		AccessRightTypeId: access_right_id,
	}

	if err := r.DB.
		FirstOrCreate(userAccessRight).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return userAccessRight, nil
}

func (r *Repo) DeleteUserAccessRight(user_id, file_id, access_right_id uint) (*model.UserAccessRight, error) {

	mux.Lock()
	defer mux.Unlock()

	userAccessRight := &model.UserAccessRight{}

	if err := r.DB.
		Where(&model.UserAccessRight{
			UserId:            user_id,
			FileId:            file_id,
			AccessRightTypeId: access_right_id,
		}).
		Delete(userAccessRight).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return userAccessRight, nil
}
