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
		GetItem(id string) (*model.Item, error)
		CreateItem(userId string, req *model.Item) (*model.Item, error)
		UpdateItem(id string, req *model.Item) (*model.Item, error)
		DeleteItem(id string) error
	}

	itemServiceImpl struct {
		repo repository.Repo
	}
)

func GetItemService() ItemService {
	if itemService == nil {
		itemService = &itemServiceImpl{
			repo: repository.GetRepository(),
		}
	}
	return itemService
}

func (is *itemServiceImpl) GetItem(id string) (*model.Item, error) {
	i, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, fmt.Errorf("not found")
	}
	return i.(*model.Item), nil
}

func (is *itemServiceImpl) CreateItem(userId string, req *model.Item) (*model.Item, error) {
	i := model.DefaultItem()
	i.UserId = &userId
	i.Label = req.Label
	i.Description = req.Description
	i.Price = req.Price
	if err := is.repo.Insert(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (is *itemServiceImpl) UpdateItem(id string, req *model.Item) (*model.Item, error) {
	obj, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return nil, err
	}
	i := obj.(*model.Item)
	// can update except soldout
	if *i.Status == model.ItemStatusSoldOut {
		return nil, fmt.Errorf("not stauts to update")
	}

	// restriction to update
	if req.Label != nil {
		i.Label = req.Label
	}
	if req.Description != nil {
		i.Description = req.Description
	}
	if req.Price != nil {
		i.Price = req.Price
	}
	if req.Status != nil {
		i.Status = req.Status
	}
	t := time.Now()
	req.UpdatedAt = &t

	if err := is.repo.Update(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (is *itemServiceImpl) DeleteItem(id string) error {
	obj, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return err
	}
	i := obj.(*model.Item)
	// can update except soldout
	if *i.Status == model.ItemStatusSoldOut {
		return fmt.Errorf("not stauts to delete")
	}

	return is.repo.Delete(i)
}
