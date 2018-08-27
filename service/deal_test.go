package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

func TestGetDealAsSeller(t *testing.T) {
	// TODO: test
}

func TestGetDealAsBuyer(t *testing.T) {
	// TODO: test
}

func TestEstablish(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()
	dealSrv := GetDealService()

	// perfect data of 1 seller has multiple items and multiple buyers
	parallel := 10
	seller, _, _, _ := createPerfectUser()
	items := make([]*model.Item, parallel)
	var pointSum int
	for idx, _ := range items {
		i := createItem(seller, model.ItemStatusSold)
		items[idx] = i
		pointSum += *i.Price
	}
	buyers := make([]*model.User, parallel)
	for idx, _ := range buyers {
		u, _, _, _ := createPerfectUser()
		buyers[idx] = u
	}

	// invalid itemId
	{
		o, e := dealSrv.Establish("hogehogeId", buyers[0].Id)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid item status
	{
		itemSoldout := createItem(seller, model.ItemStatusSoldOut)
		o, e := dealSrv.Establish(itemSoldout.Id, buyers[0].Id)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid buyterId
	{
		o, e := dealSrv.Establish(items[0].Id, "hogehogeId")
		ast.Nil(o)
		ast.NotNil(e)
	}

	// when seller's items are purchased at the same time
	// expect all deals become succeeded
	{
		ch := make(chan error, parallel)
		for idx, _ := range buyers {
			go func(i int) {
				_, err := dealSrv.Establish(items[i].Id, buyers[i].Id)
				ch <- err
			}(idx)
		}
		for range buyers {
			ast.Nil(<-ch)
		}
		as, _ := userSrv.GetAsset(seller.Id)
		ast.Equal(initialPoint+pointSum, *as.Point)
	}

	clearAll(repository.GetRepository())
}
