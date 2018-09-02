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
		SelectPublicItem(opt *Option) ([]*model.Item, error)
		FindPublicItem(id string) (*model.Item, error)
		SelectMyItem(opt *Option) ([]*model.Item, error)
		FindMyItem(id, userId string) (*model.Item, error)
	}

	itemRepoImpl struct{}
)

func GetItemRepo() ItemRepo {
	if itemRepo == nil {
		itemRepo = &itemRepoImpl{}
	}
	return itemRepo
}

func (dr *itemRepoImpl) SelectPublicItem(opt *Option) ([]*model.Item, error) {
	q := fmt.Sprintf(`
		select
		 	*
		from
			items
		where
			status != 'notsale'
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

func (ir *itemRepoImpl) FindPublicItem(id string) (*model.Item, error) {
	i := &model.Item{}
	if err := repo.dbMap.SelectOne(i, `
		select
			*
		from
			items
		where
			id = $1 and
			status != 'notsale'
		`, id); err != nil {
		return nil, err
	}
	return i, nil
}

func (dr *itemRepoImpl) SelectMyItem(opt *Option) ([]*model.Item, error) {
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

func (ir *itemRepoImpl) FindMyItem(id, userId string) (*model.Item, error) {
	i := &model.Item{}
	if err := repo.dbMap.SelectOne(i, `
		select
			*
		from
			items
		where
			id = $1 and
			user_id = $2
		`, id, userId); err != nil {
		return nil, err
	}
	return i, nil
}
