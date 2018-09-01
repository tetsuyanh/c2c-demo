package service

import (
	"sync"
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
	parallel := 3
	seller, _, _, _ := createPerfectUser()
	items := make([]*model.Item, parallel)
	for idx, _ := range items {
		i := createItem(seller, model.ItemStatusSold)
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
	{
		wg := new(sync.WaitGroup)
		ch := make(chan struct{}, parallel)
		for i := 0; i < parallel; i++ {
			wg.Add(1)
			idx := i
			go func(itemId, buyerId string) {
				defer wg.Done()
				if _, e := dealSrv.Establish(itemId, buyerId); e == nil {
					ch <- struct{}{}
				}
			}(items[idx].Id, buyers[idx].Id)
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
		cntSuccess := 0
		for range ch {
			cntSuccess++
		}
		// at least one deal become success
		ast.NotEqual(0, cntSuccess)
		ast.NotEqual(parallel, cntSuccess)
		// correct point
		as, _ := userSrv.GetAsset(seller.Id)
		ast.Equal(initialPoint+(testItemPrice*cntSuccess), as.Point)
	}

	clearAll(repository.GetRepository())
}
