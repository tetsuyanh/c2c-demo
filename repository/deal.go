package repository

import (
	"fmt"

	"github.com/tetsuyanh/c2c-demo/model"
)

const (
	DealTableName = "deals"
)

var (
	dealRepo DealRepo
)

type (
	DealRepo interface {
		SelectDealAsSeller(opt *Option) ([]*model.Deal, error)
		SelectDealAsBuyer(opt *Option) ([]*model.Deal, error)
	}

	dealRepoImpl struct{}
)

func GetDealRepo() DealRepo {
	if dealRepo == nil {
		dealRepo = &dealRepoImpl{}
	}
	return dealRepo
}

func (dr *dealRepoImpl) SelectDealAsSeller(opt *Option) ([]*model.Deal, error) {
	q := fmt.Sprintf(`
		select
		 	*
		from
			deals
		where
			seller_id = :userId
		order by
			created_at desc
		offset
			:offset
		limit
			:limit
		`)
	var ds []*model.Deal
	_, err := repo.dbMap.Select(&ds, q, opt.Map())
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func (dr *dealRepoImpl) SelectDealAsBuyer(opt *Option) ([]*model.Deal, error) {
	q := fmt.Sprintf(`
		select
		 	*
		from
			deals
		where
			buyer_id = :userId
		order by
			created_at desc
		offset
			:offset
		limit
			:limit
		`)
	var ds []*model.Deal
	_, err := repo.dbMap.Select(&ds, q, opt.Map())
	if err != nil {
		return nil, err
	}
	return ds, nil
}
