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
		GetItem(id, userId string) (*model.Item, error)
		UpdateItem(id, userId string, req *model.Item) (*model.Item, error)
		DeleteItem(id, userId string) error
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

func (is *itemServiceImpl) CreateItem(userId string, req *model.Item) (*model.Item, error) {
	i := model.DefaultItem()
	i.UserId = userId
	i.Label = req.Label
	i.Description = req.Description
	i.Price = req.Price
	if err := is.repo.Insert(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (is *itemServiceImpl) GetItem(id, userId string) (*model.Item, error) {
	obj, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return nil, err
	}
	i := obj.(*model.Item)
	if i.UserId != userId {
		return nil, fmt.Errorf("forbidden")
	}
	return i, nil
}

func (is *itemServiceImpl) UpdateItem(id, userId string, req *model.Item) (*model.Item, error) {
	obj, err := is.repo.Get(model.Item{}, id)
	if err != nil {
		return nil, err
	}
	i := obj.(*model.Item)
	if i.UserId != userId {
		return nil, fmt.Errorf("forbidden")
	}
	// updating status 'soldout' is allowed for dealService
	if i.Status == model.ItemStatusSoldOut || req.Status == model.ItemStatusSoldOut {
		return nil, fmt.Errorf("not stauts to update")
	}

	// limited fields to update
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
	req.UpdatedAt = time.Now()

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
	if i.Status == model.ItemStatusSoldOut {
		return fmt.Errorf("not stauts to delete")
	}

	return is.repo.Delete(i)
}
