package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetsuyanh/c2c-demo/model"
)

func TestGetItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()
	u, _, _, i := createPerfectUser()

	// invalid itemId
	{
		o, e := itemSrv.GetItem("hogehogeId", u.Id)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid userId
	{
		other, _, _, _ := createPerfectUser()
		o, e := itemSrv.GetItem(i.Id, other.Id)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// success
	{
		o, e := itemSrv.GetItem(i.Id, u.Id)
		ast.NotNil(o)
		ast.Nil(e)
	}
}

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
