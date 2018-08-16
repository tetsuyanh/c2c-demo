package service

import (
	"fmt"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

var (
	itemService ItemService
)

type (
	ItemService interface {
		GetItem(id string) (*model.Item, error)
		CreateItem(i *model.Item) error
		UpdateItem(i *model.Item) error
		DeleteItem(i *model.Item) error
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

func (is *itemServiceImpl) CreateItem(i *model.Item) error {
	return is.repo.Insert(i)
}

func (is *itemServiceImpl) UpdateItem(i *model.Item) error {
	return is.repo.Update(i)
}

func (is *itemServiceImpl) DeleteItem(i *model.Item) error {
	return is.repo.Delete(i)
}
