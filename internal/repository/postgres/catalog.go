package postgres

import (
	"errors"
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/pkg/errs"
)

func (r *Repo) FindCatalog(id uint) (*model.Catalog, error) {
	mux.RLock()
	defer mux.RUnlock()

	catalog := &model.Catalog{}
	if err := r.DB.
		Where("id = ?", id).
		First(catalog).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalog, nil
}

func (r *Repo) FindDiskCatalogs(disk_space_id uint) ([]*model.Catalog, error) {
	mux.RLock()
	defer mux.RUnlock()

	catalogs := []*model.Catalog{}
	if err := r.DB.
		Where(model.Catalog{DiskSpaceId:disk_space_id}).
		Find(catalogs).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalogs, nil
}

func (r *Repo) CreateCatalog(disk_space_id, parent_catalog_id uint, name string) (*model.Catalog, error) {
	mux.Lock()
	defer mux.Unlock()

	catalog := &model.Catalog{
		DiskSpaceId: disk_space_id,
		ParentCatalogId: parent_catalog_id,
		Name: name,
	}

	if err := r.DB.
		Create(catalog).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalog, nil
}

func (r *Repo) UpdateDiskCatalog(id uint, name string) (*model.Catalog, error) {
	catalog, err := r.FindCatalog(id)
	if err != nil {
		return nil, errs.NewStack(errors.New("Catalog does not exists"))
	}

	mux.Lock()
	defer mux.Unlock()

	catalog.Name = name

	if err := r.DB.
		Update(catalog).Error; err != nil {
		return nil, errs.NewStack(err)
	}

	return catalog, nil
}

func (r *Repo) DeleteDiskCatalogs(disk_space_id uint) ([]*model.Catalog, error) {
	catalogs, err := r.FindDiskCatalogs(disk_space_id)
	if err !

	mux.Lock()
	defer mux.Unlock()

	catalogs := []*model.Catalog{}
	if err := r.DB.
		Where(model.Catalog{DiskSpaceId:disk_space_id}).
		Find(catalogs).Error; err != nil {
		return nil, errs.NewStack(err)
	}
	return catalogs, nil
}
