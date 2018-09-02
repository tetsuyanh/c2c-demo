package repository

import (
	"fmt"

	"github.com/tetsuyanh/c2c-demo/model"
)

const (
	ItemTableName = "items"
)

var (
	itemRepo ItemRepo
)

type (
	ItemRepo interface {
		SelectSelfItem(opt *Option) ([]*model.Item, error)
	}

	itemRepoImpl struct{}
)

func GetItemRepo() ItemRepo {
	if itemRepo == nil {
		itemRepo = &itemRepoImpl{}
	}
	return itemRepo
}

func (dr *itemRepoImpl) SelectSelfItem(opt *Option) ([]*model.Item, error) {
	q := fmt.Sprintf(`
		select
		 	*
		from
			items
		where
			user_id = :userId
		order by
			created_at desc
		offset
			:offset
		limit
			:limit
		`)
	var ds []*model.Item
	_, err := repo.dbMap.Select(&ds, q, opt.Map())
	if err != nil {
		return nil, err
	}
	return ds, nil
}
