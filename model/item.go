package model

import (
	"fmt"
	"time"
)

var (
	ItemStatusNotSale = "notsale"
	ItemStatusSale    = "sale"
	ItemStatusSold    = "sold"
)

type Item struct {
	Id          string    `db:"id" json:"id"`
	UserId      string    `db:"user_id" json:"userId,omitempty"`
	Label       string    `db:"label" json:"label,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
	Price       int       `db:"price" json:"price,omitempty"`
	Status      string    `db:"status" json:"status,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt,omitempty"`

	// auto locking by gorp
	// see https://github.com/go-gorp/gorp#optimistic-locking
	Version int64 `db:"version" json:"-"`
}

func DefaultItem() *Item {
	t := time.Now()
	return &Item{
		generateID(),
		"",
		"",
		"",
		0,
		ItemStatusNotSale,
		t,
		t,
		0,
	}
}

func (i *Item) Verify() error {
	if i.UserId == "" {
		return fmt.Errorf("require userId")
	}
	if i.Label == "" {
		return fmt.Errorf("require label")
	}
	if i.Price <= 0 {
		return fmt.Errorf("price should be positive number")
	}
	if i.Status != ItemStatusNotSale &&
		i.Status != ItemStatusSale &&
		i.Status != ItemStatusSold {
		return fmt.Errorf("invalid status")
	}
	return nil
}
