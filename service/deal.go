package service

import (
	"fmt"

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
		userRepo repository.UserRepo
		dealRepo repository.DealRepo
	}
)

func GetDealService() DealService {
	if dealService == nil {
		dealService = &dealServiceImpl{
			repo:     repository.GetRepository(),
			userRepo: repository.GetUserRepo(),
			dealRepo: repository.GetDealRepo(),
		}
	}
	return dealService
}

func (ds *dealServiceImpl) GetDealAsSeller(opt *repository.Option) ([]*model.Deal, error) {
	if opt.GetUserId() == "" {
		return nil, fmt.Errorf("require option userId")
	}
	return ds.dealRepo.SelectDealAsSeller(opt)
}

func (ds *dealServiceImpl) GetDealAsBuyer(opt *repository.Option) ([]*model.Deal, error) {
	if opt.GetUserId() == "" {
		return nil, fmt.Errorf("require option userId")
	}
	return ds.dealRepo.SelectDealAsBuyer(opt)
}

func (ds *dealServiceImpl) Establish(itemId, buyerId string) (*model.Deal, error) {
	obj, err := ds.repo.Get(model.Item{}, itemId)
	if err != nil {
		return nil, err
	}
	tx := ds.repo.Transaction()

	i := obj.(*model.Item)
	if i.Status != model.ItemStatusSale {
		return nil, fmt.Errorf("item invalid status")
	}
	i.Status = model.ItemStatusSold
	tx.Update(i)

	d := model.DefaultDeal()
	d.ItemId = itemId
	d.SellerId = i.UserId
	d.BuyerId = buyerId
	tx.Insert(d)

	sellerAs, err := ds.userRepo.FindAssetByUserId(d.SellerId)
	if err != nil {
		return nil, err
	}
	sellerAs.Point = sellerAs.Point + i.Price
	tx.Update(sellerAs)

	buyerAs, err := ds.userRepo.FindAssetByUserId(d.BuyerId)
	if err != nil {
		return nil, err
	}
	buyerAs.Point = buyerAs.Point - i.Price
	if buyerAs.Point < 0 {
		return nil, fmt.Errorf("short of point")
	}
	tx.Update(buyerAs)

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return d, nil
}
