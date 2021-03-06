package service

import (
	"fmt"
	"time"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

var (
	itemService ItemService
)

type (
	ItemService interface {
		CreateItem(userId string, req *model.Item) (*model.Item, error)
		GetPublicItems(opt *repository.Option) ([]*model.Item, error)
		GetPublicItem(id string) (*model.Item, error)
		GetMyItems(opt *repository.Option) ([]*model.Item, error)
		GetMyItem(id, userId string) (*model.Item, error)
		UpdateItem(id, userId string, req *model.Item) (*model.Item, error)
		DeleteItem(id, userId string) error
	}

	itemServiceImpl struct {
		repo     repository.Repo
		itemRepo repository.ItemRepo
	}
)

func GetItemService() ItemService {
	if itemService == nil {
		itemService = &itemServiceImpl{
			repo:     repository.GetRepository(),
			itemRepo: repository.GetItemRepo(),
		}
	}
	return itemService
}

func (is *itemServiceImpl) CreateItem(userId string, req *model.Item) (*model.Item, error) {
	i := model.DefaultItem()
	// limited fields to create by request
	i.UserId = userId
	i.Label = req.Label
	i.Description = req.Description
	i.Price = req.Price
	if err := i.Verify(); err != nil {
		return nil, err
	}
	if err := is.repo.Insert(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (is *itemServiceImpl) GetPublicItems(opt *repository.Option) ([]*model.Item, error) {
	return is.itemRepo.SelectPublicItem(opt)
}

func (is *itemServiceImpl) GetPublicItem(id string) (*model.Item, error) {
	return is.itemRepo.FindPublicItem(id)
}

func (is *itemServiceImpl) GetMyItems(opt *repository.Option) ([]*model.Item, error) {
	if opt.GetUserId() == "" {
		return nil, fmt.Errorf("require option userId")
	}
	return is.itemRepo.SelectMyItem(opt)
}

func (is *itemServiceImpl) GetMyItem(id, userId string) (*model.Item, error) {
	i, err := is.itemRepo.FindMyItem(id, userId)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (is *itemServiceImpl) UpdateItem(id, userId string, req *model.Item) (*model.Item, error) {
	obj, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return nil, err
	}
	i := obj.(*model.Item)

	// limited fields to update by request
	if req.Label != "" {
		i.Label = req.Label
	}
	if req.Description != "" {
		i.Description = req.Description
	}
	if req.Price != 0 {
		i.Price = req.Price
	}
	if req.Status != "" {
		i.Status = req.Status
	}
	i.UpdatedAt = time.Now()
	if err := i.Verify(); err != nil {
		return nil, err
	}

	if i.UserId != userId {
		return nil, fmt.Errorf("forbidden")
	}
	// cannnot update item already sold
	// cannnot make status sold by mylsef ,is is allowed for dealService
	if i.Status == model.ItemStatusSold || req.Status == model.ItemStatusSold {
		return nil, fmt.Errorf("not stauts to update")
	}

	if err := is.repo.Update(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (is *itemServiceImpl) DeleteItem(id, userId string) error {
	obj, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return err
	}
	i := obj.(*model.Item)
	if i.UserId != userId {
		return fmt.Errorf("forbidden")
	}
	// cannnot delete item already sold
	if i.Status == model.ItemStatusSold {
		return fmt.Errorf("not stauts to delete")
	}

	return is.repo.Delete(i)
}
