package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	ast := assert.New(t)

	// invalid userId
	{
		i := perfectItem()
		i.UserId = ""
		ast.NotNil(i.Verify())
	}

	// invalid label
	{
		i := perfectItem()
		i.Label = ""
		ast.NotNil(i.Verify())
	}

	// invalid price
	{
		i := perfectItem()
		i.Price = -100
		ast.NotNil(i.Verify())
	}

	// invalid status
	{
		i := perfectItem()
		i.Status = "hoge"
		ast.NotNil(i.Verify())
	}

	// success
	{
		i := perfectItem()
		ast.Nil(i.Verify())
	}
}

func perfectItem() *Item {
	return &Item{
		UserId:      "123456780",
		Label:       "label",
		Description: "description",
		Price:       100,
		Status:      ItemStatusSale,
	}
}
