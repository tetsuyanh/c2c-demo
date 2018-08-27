package service

import (
	"fmt"
	"sync"

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
		mu       *sync.Mutex
		repo     repository.Repo
		userRepo repository.UserRepo
		dealRepo repository.DealRepo
	}
)

func GetDealService() DealService {
	if dealService == nil {
		dealService = &dealServiceImpl{
			mu:       &sync.Mutex{},
			repo:     repository.GetRepository(),
			userRepo: repository.GetUserRepo(),
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
	ds.mu.Lock()
	defer ds.mu.Unlock()

	obj, err := ds.repo.Get(model.Item{}, itemId)
	if err != nil {
		return nil, err
	}
	tx := ds.repo.Transaction()

	i := obj.(*model.Item)
	if *i.Status != model.ItemStatusSold {
		return nil, fmt.Errorf("item invalid status")
	}
	i.Status = &model.ItemStatusSoldOut
	tx.Update(i)

	d := model.DefaultDeal()
	d.ItemId = &itemId
	d.SellerId = i.UserId
	d.BuyerId = &buyerId
	tx.Insert(d)

	sellerAs, err := ds.userRepo.FindAssetByUserId(*d.SellerId)
	if err != nil {
		return nil, err
	}
	sellerPnt := *sellerAs.Point + *i.Price
	sellerAs.Point = &sellerPnt
	tx.Update(sellerAs)

	buyerAs, err := ds.userRepo.FindAssetByUserId(*d.BuyerId)
	if err != nil {
		return nil, err
	}
	buyerPnt := *buyerAs.Point - *i.Price
	buyerAs.Point = &buyerPnt
	tx.Update(buyerAs)

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return d, nil
}
