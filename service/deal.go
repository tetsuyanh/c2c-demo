package service

import (
	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

var (
	dealService DealService
)

type (
	DealService interface {
		GetDealAsSeller(opt *repository.Option) ([]*model.Deal, error)
		GetDealAsBuyer(opt *repository.Option) ([]*model.Deal, error)
		Establish(itemId, buyerId string) (*model.Deal, error)
	}

	dealServiceImpl struct {
		repo     repository.Repo
		dealRepo repository.DealRepo
	}
)

func GetDealService() DealService {
	if dealService == nil {
		dealService = &dealServiceImpl{
			repo:     repository.GetRepository(),
			dealRepo: repository.GetDealRepo(),
		}
	}
	return dealService
}

func (ds *dealServiceImpl) GetDealAsSeller(opt *repository.Option) ([]*model.Deal, error) {
	return ds.dealRepo.SelectDealAsSeller(opt)
}

func (ds *dealServiceImpl) GetDealAsBuyer(opt *repository.Option) ([]*model.Deal, error) {
	return ds.dealRepo.SelectDealAsBuyer(opt)
}

func (ds *dealServiceImpl) Establish(itemId, buyerId string) (*model.Deal, error) {
	i, err := ds.repo.Get(model.Item{}, itemId)
	if err != nil {
		return nil, err
	}
	d := model.DefaultDeal()
	d.ItemId = &itemId
	d.SellerId = i.(*model.Item).UserId
	d.BuyerId = &buyerId
	if err := ds.repo.Insert(d); err != nil {
		return nil, err
	}
	return d, nil
}
