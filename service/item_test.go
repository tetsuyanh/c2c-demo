package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

func TestCreateItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()
	u, _, _, _ := createPerfectUser()

	// invalid model
	{
		i := model.DefaultItem()
		o, e := itemSrv.CreateItem(u.Id, i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid relation
	{
		i := model.DefaultItem()
		i.Label = "label"
		i.Description = "description"
		i.Price = testItemPrice
		o, e := itemSrv.CreateItem("hogehogeId", i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// success
	{
		i := model.DefaultItem()
		i.Label = "label"
		i.Description = "description"
		i.Price = testItemPrice
		o, e := itemSrv.CreateItem(u.Id, i)
		ast.NotNil(o)
		ast.Nil(e)
	}
}

func TestGetPublicItems(t *testing.T) {
	// TODO test
	// other test functions are executed at the same time, so cannot test public items...
}

func TestGetPublicItem(t *testing.T) {
	// TODO test
	// other test functions are executed at the same time, so cannot test public items...
}

func TestGetMyItems(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()
	u := createAnonymousUser()
	createAuthentication(u.Id, true)

	// invalid userId
	{
		opt := repository.DefaultOption()
		is, e := itemSrv.GetMyItems(opt)
		ast.Nil(is)
		ast.NotNil(e)
	}

	// success empty
	{
		opt := repository.DefaultOption()
		opt.SetUserId(u.Id)
		is, e := itemSrv.GetMyItems(opt)
		ast.NotNil(is)
		ast.Equal(0, len(is))
		ast.Nil(e)
	}

	// create 3 items each status
	i1 := createItem(u.Id, model.ItemStatusNotSale)
	i2 := createItem(u.Id, model.ItemStatusSale)
	i3 := createItem(u.Id, model.ItemStatusSold)

	// success get=2, by offset=0, limit=2
	{
		opt := repository.DefaultOption()
		opt.SetUserId(u.Id)
		opt.SetOffset(0)
		opt.SetLimit(2)
		is, e := itemSrv.GetMyItems(opt)
		ast.NotNil(is)
		if ast.Equal(2, len(is)) {
			ast.Equal(i3.Id, is[0].Id)
			ast.Equal(i2.Id, is[1].Id)
		}
		ast.Nil(e)
	}

	// success get=1, by offset=2, limit=2
	{
		opt := repository.DefaultOption()
		opt.SetUserId(u.Id)
		opt.SetOffset(2)
		opt.SetLimit(2)
		is, e := itemSrv.GetMyItems(opt)
		ast.NotNil(is)
		if ast.Equal(1, len(is)) {
			ast.Equal(i1.Id, is[0].Id)
		}
		ast.Nil(e)
	}
}

func TestGetMyItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()
	u, _, _, i := createPerfectUser()

	// invalid itemId
	{
		o, e := itemSrv.GetMyItem("hogehogeId", u.Id)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid userId
	{
		other, _, _, _ := createPerfectUser()
		o, e := itemSrv.GetMyItem(i.Id, other.Id)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// success
	{
		o, e := itemSrv.GetMyItem(i.Id, u.Id)
		ast.NotNil(o)
		ast.Nil(e)
	}
}

func TestUpdateItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()
	dealSrv := GetDealService()
	u, _, _, _ := createPerfectUser()

	// invalid itemId
	{
		i := createItem(u.Id, model.ItemStatusSale)
		o, e := itemSrv.UpdateItem("hogehogeId", u.Id, i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid model
	{
		i := createItem(u.Id, model.ItemStatusSale)
		i.Price = -100
		o, e := itemSrv.UpdateItem(i.Id, u.Id, i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid userId
	{
		other := createAnonymousUser()
		i := createItem(u.Id, model.ItemStatusSale)
		o, e := itemSrv.UpdateItem(i.Id, other.Id, i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid item status
	{
		i := createItem(u.Id, model.ItemStatusSold)
		o, e := itemSrv.UpdateItem(i.Id, u.Id, i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// success
	{
		i := createItem(u.Id, model.ItemStatusSale)
		o, e := itemSrv.UpdateItem(i.Id, u.Id, i)
		ast.NotNil(o)
		ast.Nil(e)
	}

	// updating item and dealing occur at the same time
	{
		i := createItem(u.Id, model.ItemStatusSale)
		buyer, _, _, _ := createPerfectUser()
		errDeal := make(chan error)
		errUpdate := make(chan error)
		go func() {
			i.Price = 5000
			_, e := itemSrv.UpdateItem(i.Id, u.Id, i)
			errUpdate <- e
		}()
		go func() {
			_, e := dealSrv.Establish(i.Id, buyer.Id)
			errDeal <- e
		}()
		// the one should be success, the other should be fail
		if <-errDeal != nil {
			ast.Nil(<-errUpdate)
		} else {
			ast.NotNil(<-errUpdate)
		}
	}
}

func TestDeleteItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()

	u, _, _, i := createPerfectUser()

	// invalid itemId
	{
		e := itemSrv.DeleteItem("hogehogeId", u.Id)
		ast.NotNil(e)
	}

	// invalid userId
	{
		other, _, _, _ := createPerfectUser()
		e := itemSrv.DeleteItem(i.Id, other.Id)
		ast.NotNil(e)
	}

	// invalid status
	{
		itemSoldout := createItem(u.Id, model.ItemStatusSold)
		e := itemSrv.DeleteItem(itemSoldout.Id, u.Id)
		ast.NotNil(e)
	}

	// success
	{
		e := itemSrv.DeleteItem(i.Id, u.Id)
		ast.Nil(e)
	}
}
