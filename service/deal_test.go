package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

func TestGetDealAsSeller(t *testing.T) {
	ast := assert.New(t)
	dealSrv := GetDealService()
	seller, _, _, _ := createPerfectUser()
	buyer, _, _, _ := createPerfectUser()

	// invalid userId
	{
		opt := repository.DefaultOption()
		ds, e := dealSrv.GetDealAsSeller(opt)
		ast.Nil(ds)
		ast.NotNil(e)
	}

	// success empty
	{
		opt := repository.DefaultOption()
		opt.SetUserId(seller.Id)
		ds, e := dealSrv.GetDealAsSeller(opt)
		ast.NotNil(ds)
		ast.Equal(0, len(ds))
		ast.Nil(e)
	}

	// create 3 deals
	d1 := createDeal(seller.Id, buyer.Id)
	d2 := createDeal(seller.Id, buyer.Id)
	d3 := createDeal(seller.Id, buyer.Id)

	// success get=2, by offset=0, limit=2
	{
		opt := repository.DefaultOption()
		opt.SetUserId(seller.Id)
		opt.SetOffset(0)
		opt.SetLimit(2)
		ds, e := dealSrv.GetDealAsSeller(opt)
		ast.NotNil(ds)
		if ast.Equal(2, len(ds)) {
			ast.Equal(d3.Id, ds[0].Id)
			ast.Equal(d2.Id, ds[1].Id)
		}
		ast.Nil(e)
	}

	// success get=1, by offset=2, limit=2
	{
		opt := repository.DefaultOption()
		opt.SetUserId(seller.Id)
		opt.SetOffset(2)
		opt.SetLimit(2)
		ds, e := dealSrv.GetDealAsSeller(opt)
		ast.NotNil(ds)
		if ast.Equal(1, len(ds)) {
			ast.Equal(d1.Id, ds[0].Id)
		}
		ast.Nil(e)
	}
}

func TestGetDealAsBuyer(t *testing.T) {
	// TODO: test
	// almost same as TestGetDealAsSeller
}

func TestEstablish(t *testing.T) {
	ast := assert.New(t)
	// userSrv := GetUserService()
	dealSrv := GetDealService()

	// perfect data of 1 seller has multiple items and multiple buyers
	parallel := 3
	seller, _, _, _ := createPerfectUser()
	items := make([]*model.Item, parallel)
	for idx, _ := range items {
		i := createItem(seller.Id, model.ItemStatusSale)
		items[idx] = i
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
		itemSoldout := createItem(seller.Id, model.ItemStatusSold)
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
	{
		// wg := new(sync.WaitGroup)
		// ch := make(chan struct{}, parallel)
		// for i := 0; i < parallel; i++ {
		// 	wg.Add(1)
		// 	idx := i
		// 	go func(itemId, buyerId string) {
		// 		defer wg.Done()
		// 		if _, e := dealSrv.Establish(itemId, buyerId); e == nil {
		// 			ch <- struct{}{}
		// 		}
		// 	}(items[idx].Id, buyers[idx].Id)
		// }
		// go func() {
		// 	wg.Wait()
		// 	close(ch)
		// }()
		// cntSuccess := 0
		// for range ch {
		// 	cntSuccess++
		// }
		// // at least one deal become success
		// ast.NotEqual(0, cntSuccess)
		// ast.NotEqual(parallel, cntSuccess)
		// // correct point
		// as, _ := userSrv.GetAsset(seller.Id)
		// ast.Equal(initialPoint+(testItemPrice*cntSuccess), as.Point)
	}
}
