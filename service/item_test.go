package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetsuyanh/c2c-demo/model"
)

func TestGetItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()

	// invalid itemId
	{
		o, e := itemSrv.GetItem("hogehogeId")
		ast.Nil(o)
		ast.NotNil(e)
	}
	// success
	{
		_, _, _, i := createPerfectUser()
		o, e := itemSrv.GetItem(i.Id)
		ast.NotNil(o)
		ast.Nil(e)
	}
}

func TestCreateItem(t *testing.T) {
	ast := assert.New(t)
	itemSrv := GetItemService()

	u, _, _, _ := createPerfectUser()
	newItem := model.DefaultItem()
	text := "hoge"
	newItem.UserId = &u.Id
	newItem.Label = &text
	newItem.Description = &text
	newItem.Price = &testItemPrice

	// invalid user
	{

		i, e := itemSrv.CreateItem("hogehogeId", newItem)
		ast.Nil(i)
		ast.NotNil(e)
	}
	// success
	{
		i, e := itemSrv.CreateItem(u.Id, newItem)
		ast.NotNil(i)
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
		i := createItem(u, model.ItemStatusSold)
		o, e := itemSrv.UpdateItem("hogehogeId", i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// invalid item status
	{
		i := createItem(u, model.ItemStatusSoldOut)
		o, e := itemSrv.UpdateItem(i.Id, i)
		ast.Nil(o)
		ast.NotNil(e)
	}

	// success
	{
		i := createItem(u, model.ItemStatusSold)
		o, e := itemSrv.UpdateItem(i.Id, i)
		ast.NotNil(o)
		ast.Nil(e)
	}

	// updating item and dealing occur at the same time
	{
		i := createItem(u, model.ItemStatusSold)
		buyer, _, _, _ := createPerfectUser()
		errDeal := make(chan error)
		errUpdate := make(chan error)
		go func() {
			_, e := dealSrv.Establish(i.Id, buyer.Id)
			errDeal <- e
		}()
		go func() {
			price := 5000
			i.Price = &price
			_, e := itemSrv.UpdateItem(i.Id, i)
			errUpdate <- e
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
		e := itemSrv.DeleteItem("hogehogeId")
		ast.NotNil(e)
	}

	// invalid status
	{
		itemSoldout := createItem(u, model.ItemStatusSoldOut)
		e := itemSrv.DeleteItem(itemSoldout.Id)
		ast.NotNil(e)
	}

	// success
	{
		e := itemSrv.DeleteItem(i.Id)
		ast.Nil(e)
	}
}
