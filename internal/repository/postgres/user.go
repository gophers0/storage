package postgres

import (
	"errors"
	"github.com/gophers0/users/internal/model"
	"github.com/gophers0/users/pkg/errs"
	"github.com/jinzhu/gorm"
)

func (r *Repo) FindUserByLogin(login string) (*model.User, error) {
	user := &model.User{}
	if err := r.DB.Where("login = ?", login).First(user).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return user, nil
}

func (r *Repo) GetSessionTokenForUser(userId uint) (string, error) {
	user := &model.User{}
	if err := r.DB.First(user, userId).Error; err != nil {
		return "", errs.NewStack(err)
	}

	session := &model.Session{}
	err := r.DB.FirstOrCreate(session, "user_id = ?", userId).Error
	if err != nil {
		return "", errs.NewStack(err)
	}

	return session.Token, nil
}

func (r *Repo) CheckSessionTokenForUser(token string, userId uint) (*model.Session, error) {
	session := &model.Session{}
	err := r.DB.Preload("User").Where("token = ? AND user_id = ?", token, userId).First(session).Error
	if err != nil {
		return nil, errs.NewStack(err)
	}
	return session, err
}

func (r *Repo) CreateUser(login, password, role string) (*model.User, error) {
	_, err := r.FindUserByLogin(login)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// create user
		user := &model.User{
			Model:    model.Model{},
			Login:    login,
			Password: password,
			Role:     role,
		}
		if err := r.DB.Create(user).Error; err != nil {
			return nil, errs.NewStack(err)
		}
	}
	return nil, errs.NewStack(errs.UserAlreadyExists)
}
