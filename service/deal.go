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
	return ds.dealRepo.SelectDealAsSeller(opt)
}

func (ds *dealServiceImpl) GetDealAsBuyer(opt *repository.Option) ([]*model.Deal, error) {
	return ds.dealRepo.SelectDealAsBuyer(opt)
}

func (ds *dealServiceImpl) Establish(itemId, buyerId string) (*model.Deal, error) {
	obj, err := ds.repo.Get(model.Item{}, itemId)
	if err != nil {
		return nil, err
	}
	tx := ds.repo.Transaction()

	i := obj.(*model.Item)
	i.Status = &model.ItemStatusSoldOut

	d := model.DefaultDeal()
	d.ItemId = &itemId
	d.SellerId = i.UserId
	d.BuyerId = &buyerId

	sellerAs, err := ds.userRepo.FindAssetByUserId(*d.SellerId)
	if err != nil {
		return nil, err
	}
	sellerPnt := *sellerAs.Point + *i.Price
	sellerAs.Point = &sellerPnt

	buyerAs, err := ds.userRepo.FindAssetByUserId(*d.BuyerId)
	if err != nil {
		return nil, err
	}
	buyerPnt := *buyerAs.Point - *i.Price
	buyerAs.Point = &buyerPnt

	if err := tx.Update(i).Insert(d).Update(sellerAs).Update(buyerAs).Commit(); err != nil {
		return nil, err
	}
	return d, nil
}
